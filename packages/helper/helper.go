package helper

import (
	"fmt"
	"os"

	"github.com/janpreet/kado/packages/config"
)

func SetupLandingZone() error {
	landingZone := config.LandingZone
	if _, err := os.Stat(landingZone); os.IsNotExist(err) {
		err := os.MkdirAll(landingZone, 0755)
		if err != nil {
			return fmt.Errorf("failed to create landing zone: %v", err)
		}
	} else {
		err := os.RemoveAll(landingZone)
		if err != nil {
			return fmt.Errorf("failed to clean landing zone: %v", err)
		}
		err = os.MkdirAll(landingZone, 0755)
		if err != nil {
			return fmt.Errorf("failed to recreate landing zone: %v", err)
		}
	}
	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
