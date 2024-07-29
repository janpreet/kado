package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/janpreet/kado/packages/bead"
	"github.com/janpreet/kado/packages/config"
	"github.com/janpreet/kado/packages/render"
)

func isDryRun() bool {
	if len(os.Args) > 1 && os.Args[1] == "set" {
		return false
	}
	return true
}

func HandleAnsible(b bead.Bead, yamlData []map[string]interface{}, extraVarsFile bool) error {
	dryRun := isDryRun()

	playbook := b.Fields["playbook"]
	inventory := b.Fields["inventory"]
	if inventory == "" {
		inventory = filepath.Join(config.LandingZone, "inventory.ini")
	}

	args := []string{"-i", inventory}
	if extraVarsFile {
		extraVarsPath, err := render.WriteExtraVarsFile(yamlData, "yaml")
		if err != nil {
			return fmt.Errorf("failed to write extra vars file: %w", err)
		}
		args = append(args, "--extra-vars", "@"+extraVarsPath)
	}
	if dryRun {
		args = append(args, "--check")
	}
	args = append(args, filepath.Join(config.LandingZone, b.Name, playbook))

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run ansible playbook: %w", err)
	}

	return nil
}
