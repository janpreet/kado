package terragrunt

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/janpreet/kado/packages/bead"
)

func HandleTerragrunt(b bead.Bead, landingZone string, applyPlan bool) error {
    repoPath := filepath.Join(landingZone, b.Name)

    terragruntPlanPath := filepath.Join(repoPath, "plan.out")
    terragruntJSONPath := filepath.Join(repoPath, "plan.json")

    fmt.Println("Running Terragrunt plan...")
    cmd := exec.Command("terragrunt", "plan", "-out", terragruntPlanPath)
    cmd.Dir = repoPath

    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Terragrunt plan output: %s\n", string(output))
        return fmt.Errorf("failed to run Terragrunt plan: %v", err)
    }

    fmt.Println("Converting Terragrunt plan to JSON...")
    cmd = exec.Command("terragrunt", "show", "-json", terragruntPlanPath)
    cmd.Dir = repoPath

    jsonOutput, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Terragrunt show output: %s\n", string(jsonOutput))
        return fmt.Errorf("failed to convert Terragrunt plan to JSON: %v", err)
    }

    err = os.WriteFile(terragruntJSONPath, jsonOutput, 0644)
    if err != nil {
        return fmt.Errorf("failed to write JSON plan to file: %v", err)
    }

    fmt.Println("Terragrunt plan saved to:", terragruntJSONPath)

    if !applyPlan {
        fmt.Println("Skipping Terragrunt apply due to missing 'set' flag.")
        return nil
    }

    fmt.Println("Running Terragrunt apply...")
    cmd = exec.Command("terragrunt", "apply", terragruntPlanPath)
    cmd.Dir = repoPath

    applyOutput, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Printf("Terragrunt apply output: %s\n", string(applyOutput))
        return fmt.Errorf("failed to run Terragrunt apply: %v", err)
    }

    fmt.Println("Terragrunt apply completed.")
    return nil
}
