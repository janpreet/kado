package helper

import (
	"fmt"
	"os"

	"bufio"
	"log"
	"strings"
	"github.com/janpreet/kado/packages/config"
	"github.com/janpreet/kado/packages/keybase"
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

func HandleKeybaseCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: kado keybase [debug] <command>")
		fmt.Println("Available commands: link, note")
		return
	}

	
	debugIndex := -1
	for i, arg := range args {
		if arg == "debug" {
			debugIndex = i
			break
		}
	}

	if debugIndex != -1 {
		keybase.Debug = true
		
		args = append(args[:debugIndex], args[debugIndex+1:]...)
	}

	switch args[0] {
	case "link":
		err := keybase.LinkKeybase()
		if err != nil {
			log.Fatalf("Failed to link Keybase account: %v", err)
		}
		fmt.Println("Keybase account linked successfully")
	case "note":
		if len(args) < 2 {
			fmt.Println("Usage: kado keybase note <create|list|view|share>")
			return
		}
		HandleNoteCommand(args[1:])
	default:
		fmt.Printf("Unknown Keybase command: %s\n", args[0])
		fmt.Println("Available commands: link, note")
	}
}

func HandleNoteCommand(args []string) {
	switch args[0] {
	case "create":
		if len(args) < 2 {
			fmt.Println("Usage: kado keybase note create <note_name>")
			return
		}
		noteName := args[1]
		fmt.Println("Enter note content (press Ctrl+D when finished):")
		scanner := bufio.NewScanner(os.Stdin)
		var content strings.Builder
		for scanner.Scan() {
			content.WriteString(scanner.Text() + "\n")
		}
		err := keybase.CreateNote(noteName, content.String())
		if err != nil {
			log.Fatalf("Failed to create note: %v", err)
		}
		fmt.Println("Note created successfully")

	case "list":
		notes, err := keybase.ListNotes()
		if err != nil {
			log.Fatalf("Failed to list notes: %v", err)
		}
		if len(notes) == 0 {
			fmt.Println("No notes found")
		} else {
			fmt.Println("Stored notes:")
			for _, note := range notes {
				fmt.Println("-", note)
			}
		}

	case "view":
		if len(args) < 2 {
			fmt.Println("Usage: kado keybase note view <note_name>")
			return
		}
		noteName := args[1]
		content, err := keybase.ViewNote(noteName)
		if err != nil {
			log.Fatalf("Failed to view note: %v", err)
		}
		fmt.Printf("Content of note '%s':\n%s", noteName, content)

	case "share":
		if len(args) < 3 {
			fmt.Println("Usage: kado keybase note share <note_name> <keybase_username>")
			return
		}
		noteName := args[1]
		recipient := args[2]
		err := keybase.ShareNote(noteName, recipient)
		if err != nil {
			log.Fatalf("Failed to share note: %v", err)
		}
		fmt.Printf("Note '%s' shared with %s successfully\n", noteName, recipient)

    case "create-with-tags":
        if len(args) < 3 {
            fmt.Println("Usage: kado keybase note create-with-tags <note_name> <tag1,tag2,...>")
            return
        }
        noteName := args[1]
        tags := strings.Split(args[2], ",")
        fmt.Println("Enter note content (press Ctrl+D when finished):")
        content := readMultiLineInput()
        err := keybase.CreateNoteWithTags(noteName, keybase.Note{Content: content, Tags: tags})
        if err != nil {
            log.Fatalf("Failed to create note with tags: %v", err)
        }
        fmt.Println("Note created successfully with tags")
    case "search-by-tag":
        if len(args) < 2 {
            fmt.Println("Usage: kado keybase note search-by-tag <tag>")
            return
        }
        tag := args[1]
        notes, err := keybase.SearchNotesByTag(tag)
        if err != nil {
            log.Fatalf("Failed to search notes by tag: %v", err)
        }
        fmt.Printf("Notes with tag '%s':\n", tag)
        for _, note := range notes {
            fmt.Printf("  - %s\n", note)
        }		

	default:
		fmt.Printf("Unknown note command: %s\n", args[0])
	}
}

func readMultiLineInput() string {
    var content strings.Builder
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        content.WriteString(scanner.Text() + "\n")
    }
    return content.String()
}