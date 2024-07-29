package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/janpreet/kado/packages/bead"
)

func HandleTerraform(b bead.Bead, landingZone string, applyPlan bool) error {
	fmt.Printf("Processing terraform bead:\n")
	for key, val := range b.Fields {
		fmt.Printf("  %s = %s\n", key, val)
	}

	fmt.Println("Getting tfvars files from landing zone:", landingZone)
	varFiles, err := getTfvarsFiles(landingZone)
	if err != nil {
		return fmt.Errorf("failed to get tfvars files: %v", err)
	}

	repoPath := filepath.Join(landingZone, b.Name)
	planArgs := []string{"plan", "-out=plan.out"}
	for _, varFile := range varFiles {

		destPath := filepath.Join(repoPath, filepath.Base(varFile))
		err := moveFile(varFile, destPath)
		if err != nil {
			return fmt.Errorf("failed to move tfvars file: %v", err)
		}
		planArgs = append(planArgs, "--var-file", filepath.Base(varFile))
	}

	backendConfigFile := filepath.Join(landingZone, "backend.tfvars")
	if fileExists(backendConfigFile) {
		destBackendPath := filepath.Join(repoPath, "backend.tfvars")
		err := moveFile(backendConfigFile, destBackendPath)
		if err != nil {
			return fmt.Errorf("failed to move backend.tfvars file: %v", err)
		}
	}

	initArgs := []string{"init"}

	if fileExists(filepath.Join(repoPath, "backend.tfvars")) {
		initArgs = append(initArgs, "-backend-config=backend.tfvars")
	}

	fmt.Println("Running terraform init...")
	err = runCommand(repoPath, "terraform", initArgs...)
	if err != nil {
		return fmt.Errorf("failed to run terraform init: %v", err)
	}

	fmt.Println("Running terraform plan...")
	err = runCommand(repoPath, "terraform", planArgs...)
	if err != nil {
		return fmt.Errorf("failed to run terraform plan: %v", err)
	}

	fmt.Println("Converting plan.out to plan.json...")
	showArgs := []string{"show", "-no-color", "-json", "plan.out"}
	output, err := runCommandWithOutput(repoPath, "terraform", showArgs...)
	if err != nil {
		return fmt.Errorf("failed to run terraform show: %v", err)
	}

	planJSONPath := filepath.Join(repoPath, "plan.json")
	err = os.WriteFile(planJSONPath, output, 0644)
	if err != nil {
		return fmt.Errorf("failed to write plan.json: %v", err)
	}

	fmt.Println("Terraform plan saved as plan.json")

	if applyPlan {

		applyArgs := []string{"apply", "plan.out"}
		fmt.Println("Applying terraform plan...")
		err = runCommand(repoPath, "terraform", applyArgs...)
		if err != nil {
			return fmt.Errorf("failed to apply terraform plan: %v", err)
		}
	}

	return nil
}

func runCommand(dir, name string, args ...string) error {
	fmt.Printf("Executing command: %s %s in directory: %s\n", name, strings.Join(args, " "), dir)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandWithOutput(dir, name string, args ...string) ([]byte, error) {
	fmt.Printf("Executing command: %s %s in directory: %s\n", name, strings.Join(args, " "), dir)
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
}

func getTfvarsFiles(directory string) ([]string, error) {
	fmt.Println("Reading tfvars files from directory:", directory)
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var varFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".tfvars") && file.Name() != "backend.tfvars" {
			varFiles = append(varFiles, filepath.Join(directory, file.Name()))
		}
	}
	fmt.Println("Found tfvars files:", varFiles)
	return varFiles, nil
}

func moveFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}
	err = os.Remove(src)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
