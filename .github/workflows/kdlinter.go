package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type LintError struct {
	File    string
	Line    int
	Message string
}

func (e LintError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}

func LintKDFile(filePath string) ([]LintError, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lintErrors []LintError
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	insideBead := false
	commentRegex := regexp.MustCompile(`^\s*#`)
	emptyLineCount := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "bead ") && strings.HasSuffix(trimmedLine, "{") {
			if insideBead {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Nested beads are not allowed"})
			}
			insideBead = true
			if !strings.HasPrefix(line, "bead ") {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Bead declaration should start at the beginning of the line"})
			}
		} else if trimmedLine == "}" && insideBead {
			insideBead = false
			if !strings.HasPrefix(line, "}") {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Closing brace should be at the beginning of the line"})
			}
		} else if insideBead {
			if !strings.HasPrefix(line, "  ") {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Bead content should be indented with two spaces"})
			}
		} else {
			if commentRegex.MatchString(trimmedLine) {
				if !strings.HasPrefix(trimmedLine, "#") {
					lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Comments should start at the beginning of the line"})
				}
			} else if trimmedLine != "" && strings.HasPrefix(line, " ") {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Non-bead content should not be indented"})
			}
		}

		if trimmedLine == "" {
			emptyLineCount++
			if emptyLineCount > 1 {
				lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Multiple consecutive empty lines are not allowed"})
			}
		} else {
			emptyLineCount = 0
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if insideBead {
		lintErrors = append(lintErrors, LintError{filePath, lineNumber, "Unclosed bead at end of file"})
	}

	return lintErrors, nil
}

func LintKDFilesInDir(dir string) ([]LintError, error) {
	var allErrors []LintError

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".kd") {
			errors, err := LintKDFile(path)
			if err != nil {
				return fmt.Errorf("failed to lint %s: %v", path, err)
			}
			allErrors = append(allErrors, errors...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allErrors, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: kdlinter <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	errors, err := LintKDFilesInDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	for _, e := range errors {
		fmt.Println(e)
	}

	if len(errors) > 0 {
		os.Exit(1)
	}
}