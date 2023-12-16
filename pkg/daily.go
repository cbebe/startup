package startup

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"slices"
)

func dailyNotesDir() (string, error) {
	h := os.Getenv("HOME")
	if h == "" {
		return "", fmt.Errorf("HOME is not defined")
	}
	return path.Join(h, "notes", "daily"), nil
}

func getMostRecentEntry() (*Entry, error) {
	s, err := getDailyNotes()
	if err != nil {
		return nil, fmt.Errorf("getDailyNotes: %v", err)
	}
	recent := slices.Max(s)
	d, err := dailyNotesDir()
	if err != nil {
		return nil, fmt.Errorf("dailyNotesDir: %v", err)
	}
	filePath := path.Join(d, recent)
	e, err := entryFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("entryFromFile: %v", err)
	}
	return e, nil
}

func getDailyNotes() ([]string, error) {
	d, err := dailyNotesDir()
	if err != nil {
		return nil, fmt.Errorf("dailyNotesDir: %v", err)
	}
	// Only get daily notes
	matches, err := fs.Glob(os.DirFS(d), "*-*-*.md")
	if err != nil {
		return nil, fmt.Errorf("fs.Glob: %v", err)
	}
	return matches, nil
}

func execCommand(command string) error {
	cmd := exec.Command("zsh")
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdinPipe: %v", err)
	}
	_, err = in.Write([]byte(command + "\n"))
	if err != nil {
		return fmt.Errorf("write to pipe: %v", err)
	}
	err = in.Close()
	if err != nil {
		return fmt.Errorf("error closing stdin: %v", err)
	}
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running command: %v", err)
	}
	return nil
}
