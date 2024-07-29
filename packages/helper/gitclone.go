package helper

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneRepo(source, destination, beadName, refs string) error {

	beadDir := filepath.Join(destination, beadName)
	if !FileExists(beadDir) {
		err := os.MkdirAll(beadDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "clone", source, beadDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if refs != "" {
		cmd = exec.Command("git", "-C", beadDir, "checkout", refs)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	log.Printf("Repository for %s cloned to: %s", beadName, beadDir)
	return nil
}
