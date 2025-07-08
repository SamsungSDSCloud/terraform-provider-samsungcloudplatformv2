package test

import (
	"bufio"
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/internal/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type testResult int

const (
	TestFailed testResult = iota
	TestPassed
)

type TestCases struct {
	Create Create `yaml:"create"`
	Update Update `yaml:"update"`
}

// Create
type Create struct {
	Unhappy []string `yaml:"unhappy"`
	Happy   string   `yaml:"happy"`
}

// Update
type Update struct {
	Unhappy []string `yaml:"unhappy"`
	Happy   []string `yaml:"happy"`
}

func getTestCases(filepath string) (*TestCases, error) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	tc := &TestCases{}
	err = yaml.Unmarshal(buf, tc)
	if err != nil {
		return nil, err
	}

	return tc, nil
}

var useStep = false
var passFailString = [2]string{"Fail", "Pass"}

var resultCnt = [2]int{0, 0}
var totalResultCnt = [2]int{0, 0}

func writeResult(w *bufio.Writer, str string, result testResult) {
	w.WriteString(str + " : " + passFailString[result] + "\n")
	resultCnt[result]++
}

func runTestList(testCases []string, testTitle string, absDir string, pass testResult, fail testResult, w *bufio.Writer) {

	kb := bufio.NewReader(os.Stdin)

	for _, tc := range testCases {
		varFile := "-var-file=" + absDir + string(os.PathSeparator) + "tfvars" + string(os.PathSeparator) + tc
		err := RunTerraformCommand([]string{"apply", "--auto-approve", varFile}, absDir, true)

		tcName := testTitle + " : " + tc[0:len(tc)-len(".tfvars")]
		if err == nil {
			writeResult(w, tcName, pass)
		} else {
			writeResult(w, tcName, fail)
		}

		if useStep {
			fmt.Printf("%s done...\n", tc)
			fmt.Println("Press enter to continue...")
			_, err = kb.ReadString('\n')
			if err != nil {
				return
			}
		}
	}
}

func drawLine(width int, c string) {
	for i := 0; i < width; i++ {
		fmt.Printf(c)
	}

	fmt.Printf("\n")
}

func PrintTestList(root string) error {
	var targetList []string
	maxWidth := 10

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if file == "testcases.yaml" {
			targetName := dir[0 : len(dir)-1]
			targetList = append(targetList, targetName)
			maxWidth = util.Max(maxWidth, len(targetName))
		}
		return nil
	})

	fmt.Printf("Testing targets\n")
	drawLine(maxWidth+2, "=")
	for _, t := range targetList {
		fmt.Printf(t + "\n")
	}

	return nil
}

func RunResourceTest(root string, targets []string, step bool) error {
	var resultFile *os.File
	var err error = nil
	defer func() {
		if err != nil {
			fmt.Errorf(err.Error())
		}
		resultFile.Close()
	}()

	now := time.Now().Format(time.RFC3339)
	now = strings.Replace(now, ":", "_", -1)

	testFileName := fmt.Sprintf("ResultSummary_%s.txt", now)
	_, err = os.Stat(testFileName)
	if os.IsNotExist(err) == false {
		os.Remove(testFileName)
	}

	resultFile, err = os.OpenFile(testFileName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.FileMode(644))
	if err != nil {
		fmt.Println("fail to open file:", err)
		return err
	}

	useStep = step
	w := bufio.NewWriter(resultFile)

	kb := bufio.NewReader(os.Stdin)

	testTargets := make(map[string]struct{})

	if len(targets) > 0 {
		for _, target := range targets {
			testTargets[target] = struct{}{}
		}
	}

	maxLength := 10
	var testResults []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if file != "testcases.yaml" {
			return nil
		}

		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil
		}

		if len(testTargets) > 0 {
			curProd := filepath.Base(dir)
			_, ok := testTargets[curProd]
			if ok == false {
				return nil
			}
		}

		err = RemoveByproducts(absDir)
		if err != nil {
			fmt.Println("fail to clear directory:", err)
			return err
		}

		err = RunTerraformCommand([]string{"init"}, absDir, true)
		if err != nil {
			fmt.Println("fail to init terraform:", err)
			return err
		}

		tcs, err := getTestCases(path)
		if err != nil {
			fmt.Println("fail to get test cases from "+path+":", err)
			return err
		}

		resultCnt[0] = 0
		resultCnt[1] = 0

		testName := dir[0 : len(dir)-1]
		w.WriteString("[" + testName + " Tests" + "]" + "\n")

		// create fail
		runTestList(tcs.Create.Unhappy, "Create fail test", absDir, TestFailed, TestPassed, w)

		// create
		varFile := "-var-file=" + absDir + string(os.PathSeparator) + "tfvars" + string(os.PathSeparator) + tcs.Create.Happy
		err = RunTerraformCommand([]string{"apply", "--auto-approve", varFile}, absDir, true)
		if err != nil {
			writeResult(w, "Create Pass test", TestFailed)
			return err
		} else {
			writeResult(w, "Create Pass test", TestPassed)
		}

		if useStep {
			fmt.Printf("%s done...\n", tcs.Create.Happy)
			fmt.Println("Press enter to continue...")
			_, inErr := kb.ReadString('\n')
			if inErr != nil {
				return err
			}
		}

		// update fail
		runTestList(tcs.Update.Unhappy, "Update fail test", absDir, TestFailed, TestPassed, w)

		// update pass
		runTestList(tcs.Update.Happy, "Update pass test", absDir, TestPassed, TestFailed, w)

		// destroy
		varFile = "-var-file=" + absDir + string(os.PathSeparator) + "tfvars" + string(os.PathSeparator) + "destroy.tfvars"
		err = RunTerraformCommand([]string{"destroy", "--auto-approve", varFile}, absDir, true)
		if err != nil {
			writeResult(w, "destroy test", TestFailed)
			return err
		} else {
			writeResult(w, "destroy test", TestPassed)
		}

		err = RemoveByproducts(absDir)
		if err != nil {
			fmt.Println("fail to clear directory:", err)
			return err
		}

		if useStep {
			fmt.Println("destroy.tfvars done...")
			fmt.Println("Press enter to continue...")
			_, inErr := kb.ReadString('\n')
			if inErr != nil {
				return err
			}
		}

		totalResultCnt[0] += resultCnt[0]
		totalResultCnt[1] += resultCnt[1]

		summary := fmt.Sprintf("["+testName+" result] : %d case(s) passed, /%d case(s) failed\n", resultCnt[1], resultCnt[0])
		testResults = append(testResults, summary)

		maxLength = util.Max(maxLength, len(summary))
		w.WriteString(summary)

		return nil
	})

	summary := fmt.Sprintf("[Total Summary]  : %d case(s) passed/%d case(s) failed\n", totalResultCnt[1], totalResultCnt[0])
	w.WriteString("--------------------------------------------------\n")
	w.WriteString("--------------------------------------------------\n")
	w.WriteString(summary)
	w.Flush()

	resultFile.Close()

	if len(testResults) > 0 {
		drawLine(maxLength+2, "=")
		for _, t := range testResults {
			fmt.Printf(t)
		}
		fmt.Printf("[Total Summary]  : %d case(s) passed/%d case(s) failed\n", totalResultCnt[1], totalResultCnt[0])
	}

	return nil
}
