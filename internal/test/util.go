package test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func RemoveByproducts(dir string) error {
	dir += string(os.PathSeparator)
	targets := []string{
		dir + ".terraform.lock.hcl",
		dir + "terraform.tfstate",
		dir + "terraform.lock.info",
		dir + ".terraform",
		dir + "terraform.tfstate.backup",
	}

	for _, target := range targets {
		_, err := os.Stat(target)
		if os.IsNotExist(err) == false {
			fmt.Printf("Removing %s...\n", target)
			err = os.RemoveAll(target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RunTerraformCommand(args []string, dir string, print bool) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	if print {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = io.Discard
	}

	return cmd.Run()
}
