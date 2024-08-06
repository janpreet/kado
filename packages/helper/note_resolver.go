package helper

import (
    "regexp"
    "strings"
	"fmt"
    "github.com/janpreet/kado/packages/keybase"
)

var noteReferenceRegex = regexp.MustCompile(`{{keybase:note:([^}]+)}}`)

func resolveNoteReferences(input string) (string, error) {
    return noteReferenceRegex.ReplaceAllStringFunc(input, func(match string) string {
        noteName := strings.TrimPrefix(strings.TrimSuffix(match, "}}"), "{{keybase:note:")
        content, err := keybase.ViewNote(noteName)
        if err != nil {
            
            return fmt.Sprintf("ERROR: Could not resolve note %s", noteName)
        }
        return strings.TrimSpace(content)
    }), nil
}