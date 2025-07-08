package testcli

import (
	"fmt"
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatformv2/internal/test"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "run-testcases",
	Short: "Command line interface for testing SCP terraform provider",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run-testcases")
		fmt.Println("Use \"run-testcases --help for more details\"")
	},
}

func NewCmdRunDatasourceTest() *cobra.Command {
	var root string
	var targets []string
	var step string

	cmd := &cobra.Command{
		Use:   "datasource",
		Short: "Run all data-sources in the given directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			useStep := false
			if strings.EqualFold(step, "yes") {
				useStep = true
			}
			err := test.RunDatasourceTest(root, targets, useStep)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&root, "root", "r", "../test", "Root directory containing test cases")
	cmd.Flags().StringSliceVarP(&targets, "target services", "t", []string{}, "target services for test")
	cmd.Flags().StringVarP(&step, "step", "s", "no", "Step-by-step testing")

	cmd.MarkFlagRequired("root")

	return cmd
}

func NewCmdRunResourceTest() *cobra.Command {
	var root string
	var targets []string
	var step string

	cmd := &cobra.Command{
		Use:   "resource",
		Short: "Run all resource testcases in the given directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			useStep := false
			if strings.EqualFold(step, "yes") {
				useStep = true
			}
			err := test.RunResourceTest(root, targets, useStep)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&root, "root", "r", "../test", "Root directory containing test cases")
	cmd.Flags().StringSliceVarP(&targets, "target services", "t", []string{}, "target services for test")
	cmd.Flags().StringVarP(&step, "step", "s", "no", "Step-by-step testing")

	cmd.MarkFlagRequired("root")

	return cmd
}

func NewCmdTestList() *cobra.Command {
	var root string

	cmd := &cobra.Command{
		Use:   "targets",
		Short: "List testcases in the given directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			err := test.PrintTestList(root)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&root, "root", "r", "../test", "Root directory containing test cases")
	cmd.MarkFlagRequired("root")

	return cmd
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewCmdRunResourceTest())
	rootCmd.AddCommand(NewCmdTestList())
	rootCmd.AddCommand(NewCmdRunDatasourceTest())
}
