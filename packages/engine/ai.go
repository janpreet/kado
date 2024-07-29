package engine

import (
    "fmt"
    "log"

    kadoai "github.com/janpreet/kado-ai/ai"
	"github.com/janpreet/kado/packages/config"
)

func RunAI() {
    client, err := kadoai.NewAIClient(config.LandingZone, "")
    if err != nil {
        log.Fatalf("Error creating AI client: %v", err)
    }

    recommendations, err := client.RunAI()
    if err != nil {
        if err.Error() == "operation cancelled by user" {
            fmt.Println("AI analysis cancelled.")
            return
        }
        log.Fatalf("Error running AI: %v", err)
    }

    fmt.Println("Infrastructure Recommendations:")
    fmt.Println(recommendations)
}