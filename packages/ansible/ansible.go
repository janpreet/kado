package ansible

import (
	"fmt"
	"os/exec"
	"strings"
)

func RunPlaybook(playbookPath, inventoryPath, extraVarsPath string, dryRun bool) error {
	args := []string{playbookPath}
	if inventoryPath != "" {
		args = append(args, "-i", inventoryPath)
	}
	args = append(args, "--extra-vars", fmt.Sprintf("@%s", extraVarsPath))
	if dryRun {
		args = append(args, "--check")
	}

	cmd := exec.Command("ansible-playbook", args...)
	fmt.Printf("Running Ansible command: ansible-playbook %s\n", strings.Join(args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run ansible playbook: %v, output: %s", err, string(output))
	}

	fmt.Printf("Ansible playbook completed. Output:\n%s\n", string(output))
	return nil
}
