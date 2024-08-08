package helper

import (
    "regexp"
    "strings"
	"fmt"
	"log"
    "github.com/janpreet/kado/packages/keybase"
	"github.com/janpreet/kado/packages/config"
)

var noteReferenceRegex = regexp.MustCompile(`{{keybase:note:([^}]+)}}`)

func resolveKeybaseNote(noteName string) (string, error) {
    if config.Debug {
        log.Printf("Attempting to resolve Keybase note: %s", noteName)
    }

    content, err := keybase.ViewNote(noteName)
    if err != nil {
        if kerr, ok := err.(*keybase.KeybaseError); ok {
            switch kerr.Type {
            case keybase.ErrNoteNotFound:
                log.Printf("WARNING: Keybase note '%s' not found", noteName)
                return "", fmt.Errorf("Keybase note '%s' not found", noteName)
            case keybase.ErrKeybaseNotInitialized:
                log.Printf("ERROR: Keybase is not initialized. Please run 'kado keybase link' first")
                return "", fmt.Errorf("Keybase is not initialized")
            case keybase.ErrPermissionDenied:
                log.Printf("ERROR: Permission denied when accessing Keybase note '%s'", noteName)
                return "", fmt.Errorf("Permission denied for Keybase note '%s'", noteName)
            default:
                log.Printf("ERROR: Failed to retrieve Keybase note '%s': %v", noteName, err)
                return "", fmt.Errorf("Failed to retrieve Keybase note '%s': %v", noteName, err)
            }
        } else {
            log.Printf("ERROR: Unexpected error when resolving Keybase note '%s': %v", noteName, err)
            return "", fmt.Errorf("Unexpected error when resolving Keybase note '%s': %v", noteName, err)
        }
    }

    if config.Debug {
        log.Printf("Successfully resolved Keybase note: %s", noteName)
    }

    return strings.TrimSpace(content), nil
}