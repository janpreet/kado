package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FormatKDFile formats a single .kd file with proper indentation and without extra newlines at the end.
func FormatKDFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var formattedLines []string
	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	insideBead := false
	commentRegex := regexp.MustCompile(`^\s*#`)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "bead ") && strings.HasSuffix(trimmedLine, "{") {
			if insideBead {
				buffer.WriteString("}\n")
				formattedLines = append(formattedLines, buffer.String())
				buffer.Reset()
			}
			insideBead = true
			buffer.WriteString(trimmedLine + "\n")
		} else if trimmedLine == "}" && insideBead {
			buffer.WriteString(trimmedLine + "\n")
			formattedLines = append(formattedLines, buffer.String())
			buffer.Reset()
			insideBead = false
		} else if insideBead {
			buffer.WriteString("  " + trimmedLine + "\n")
		} else {
			if commentRegex.MatchString(trimmedLine) {
				// Remove extra spaces in comments
				trimmedLine = strings.TrimSpace(trimmedLine)
				trimmedLine = strings.Replace(trimmedLine, "#", "#", 1)
			}
			if trimmedLine != "" || (len(formattedLines) > 0 && formattedLines[len(formattedLines)-1] != "\n") {
				formattedLines = append(formattedLines, trimmedLine+"\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if insideBead {
		buffer.WriteString("}\n")
		formattedLines = append(formattedLines, buffer.String())
	}

	// Remove the last newline if it exists
	if len(formattedLines) > 0 && formattedLines[len(formattedLines)-1] == "\n" {
		formattedLines = formattedLines[:len(formattedLines)-1]
	}

	formattedContent := strings.Join(formattedLines, "")

	err = os.WriteFile(filePath, []byte(formattedContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// FormatKDFilesInDir formats all .kd files in a given directory with proper indentation and without extra newlines at the end.
func FormatKDFilesInDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".kd") {
			err := FormatKDFile(path)
			if err != nil {
				return fmt.Errorf("failed to format %s: %v", path, err)
			}
			fmt.Printf("Formatted: %s\n", path)
		}
		return nil
	})
}
