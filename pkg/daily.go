package startup

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/bitfield/script"
)

func getDailyNotesPipe() *script.Pipe {
	h := os.Getenv("HOME")
	p := script.ListFiles(path.Join(h, "notes", "daily", "*-*-*.md"))
	if h == "" {
		p.SetError(fmt.Errorf("HOME is not defined"))
	}
	return p
}

func getJournalPrompts() ([]string, error) {
	s, err := script.File("journal.txt").Reject("#").String()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(s), "\n"), nil
}

func getMostRecentEntry() (*Entry, error) {
	recent, err := getDailyNotesPipe().Last(1).String()
	if err != nil {
		return nil, fmt.Errorf("getDailyNotesPipe: %v", err)
	}
	if recent == "" {
		return nil, fmt.Errorf("glob not found")
	}
	e, err := entryFromFile(recent)
	if err != nil {
		return nil, fmt.Errorf("entryFromFile: %v", err)
	}
	return e, nil
}

func getDailyNotes() ([]string, error) {
	f, err := getDailyNotesPipe().String()
	if err != nil {
		return nil, fmt.Errorf("getDailyNotesPipe: %v", err)
	}
	return strings.Split(strings.TrimSpace(f), "\n"), nil
}
