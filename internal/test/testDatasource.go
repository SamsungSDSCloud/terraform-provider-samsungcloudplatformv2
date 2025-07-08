package test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func RunDatasourceTest(root string, targets []string, step bool) error {

	testTargets := make(map[string]struct{})

	if len(targets) > 0 {
		for _, target := range targets {
			testTargets[target] = struct{}{}
		}
	}

	kb := bufio.NewReader(os.Stdin)
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if file != "main.tf" {
			return nil
		}

		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil
		}

		curProd := filepath.Base(dir)

		if len(testTargets) > 0 {
			_, ok := testTargets[curProd]
			if ok == false {
				return nil
			}
		}
		fmt.Printf("Prepare %s...\n", curProd)

		err = RemoveByproducts(absDir)
		if err != nil {
			fmt.Println("fail to clear directory:", err)
			return err
		}

		err = RunTerraformCommand([]string{"init"}, absDir, false)
		if err != nil {
			fmt.Println("fail to init terraform:", err)
			return err
		}

		varFile := "-var-file=" + absDir + string(os.PathSeparator) + "console.tfvars"
		err = RunTerraformCommand([]string{"apply", "--auto-approve", varFile}, absDir, false)
		if err != nil {
			fmt.Println("fail to apply datasource:", err)
			return err
		}

		fmt.Printf("Applying %s done...\n", curProd)

		if step {
			fmt.Println("Press enter to continue...")
			_, inErr := kb.ReadString('\n')
			if inErr != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
