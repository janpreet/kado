package keybase

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Debug bool

type Note struct {
    Content string
    Tags    []string
}

func CreateNoteWithTags(noteName string, note Note) error {
    content := fmt.Sprintf("Tags: %s\n\n%s", strings.Join(note.Tags, ", "), note.Content)
    return CreateNote(noteName, content)
}

func GetNoteTags(noteName string) ([]string, error) {
    content, err := ViewNote(noteName)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(content, "\n")
    if len(lines) > 0 && strings.HasPrefix(lines[0], "Tags: ") {
        tags := strings.TrimPrefix(lines[0], "Tags: ")
        return strings.Split(tags, ", "), nil
    }

    return []string{}, nil
}

func SearchNotesByTag(tag string) ([]string, error) {
    notes, err := ListNotes()
    if err != nil {
        return nil, err
    }

    var matchingNotes []string
    for _, note := range notes {
        tags, err := GetNoteTags(note)
        if err != nil {
            return nil, err
        }
        for _, t := range tags {
            if t == tag {
                matchingNotes = append(matchingNotes, note)
                break
            }
        }
    }

    return matchingNotes, nil
}

func LinkKeybase() error {
    cmd := exec.Command("keybase", "login")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("failed to link Keybase account: %v", err)
    }
    return nil
}

func InitNoteRepository() error {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %v", err)
    }
    notesDir := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes")
    
    cmd := exec.Command("git", "init")
    cmd.Dir = notesDir
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to initialize git repository: %v", err)
    }
    return nil
}

func CheckKeybaseSetup() error {
	cmd := exec.Command("keybase", "status")
	output, err := cmd.CombinedOutput()
	if Debug {
		fmt.Printf("Keybase status output:\n%s\n", string(output))
	}
	if err != nil {
		return fmt.Errorf("Keybase is not properly set up: %v", err)
	}
	if !strings.Contains(string(output), "Logged in:     yes") {
		return fmt.Errorf("You are not logged in to Keybase. Please run 'kado keybase link' first")
	}
	return nil
}

func CreateNote(noteName, content string) error {
	if err := CheckKeybaseSetup(); err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	notePath := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes", noteName)
	dir := filepath.Dir(notePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
    if err := os.WriteFile(notePath, []byte(content), 0600); err != nil {
        return fmt.Errorf("failed to write note: %v", err)
    }

    
    if err := gitAddCommit(notePath, "Create note "+noteName); err != nil {
        return fmt.Errorf("failed to version note: %v", err)
    }

    if Debug {
        fmt.Printf("Note created and versioned at: %s\n", notePath)
    }
    return nil
}

func ListNotes() ([]string, error) {
	if err := CheckKeybaseSetup(); err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}
	notesDir := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes")
	if _, err := os.Stat(notesDir); os.IsNotExist(err) {
		return []string{}, nil
	}
	files, err := os.ReadDir(notesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes directory: %v", err)
	}
	var notes []string
	for _, file := range files {
		if !file.IsDir() {
			notes = append(notes, file.Name())
		}
	}
	return notes, nil
}

func ViewNote(noteName string) (string, error) {
	if err := CheckKeybaseSetup(); err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	notePath := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes", noteName)
	content, err := os.ReadFile(notePath)
	if err != nil {
		return "", fmt.Errorf("failed to read note: %v", err)
	}
	return string(content), nil
}

func ShareNote(noteName, recipient string) error {
	if err := CheckKeybaseSetup(); err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	sourcePath := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes", noteName)
	destPath := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), recipient, "kado_notes", noteName)
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}
	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source note: %v", err)
	}
	err = os.WriteFile(destPath, input, 0600)
	if err != nil {
		return fmt.Errorf("failed to write shared note: %v", err)
	}
	if Debug {
		fmt.Printf("Note shared at: %s\n", destPath)
	}
	return nil
}

func UpdateNote(noteName, content string) error {
	if err := CheckKeybaseSetup(); err != nil {
		return err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	notePath := filepath.Join(homeDir, "Keybase", "private", os.Getenv("USER"), "kado_notes", noteName)
	dir := filepath.Dir(notePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
    if err := os.WriteFile(notePath, []byte(content), 0600); err != nil {
        return fmt.Errorf("failed to write note: %v", err)
    }

    if err := gitAddCommit(notePath, "Update note "+noteName); err != nil {
        return fmt.Errorf("failed to version note update: %v", err)
    }
    return nil
}

func gitAddCommit(filePath, message string) error {
    dir := filepath.Dir(filePath)
    
    addCmd := exec.Command("git", "add", filepath.Base(filePath))
    addCmd.Dir = dir
    if err := addCmd.Run(); err != nil {
        return fmt.Errorf("git add failed: %v", err)
    }

    commitCmd := exec.Command("git", "commit", "-m", message)
    commitCmd.Dir = dir
    if err := commitCmd.Run(); err != nil {
        return fmt.Errorf("git commit failed: %v", err)
    }

    return nil
}