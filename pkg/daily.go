package startup

import (
	"fmt"
	"io/fs"
	"log"
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

func LinkMain() {
	s, err := getDailyNotes()
	if err != nil {
		log.Fatal(err)
	}
	recent := slices.Max(s)
	d, err := dailyNotesDir()
	if err != nil {
		log.Fatal(err)
	}
	filePath := path.Join(d, recent)
	e, err := entryFromFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e.FrontMatter())
}

func getDailyNotes() ([]string, error) {
	d, err := dailyNotesDir()
	if err != nil {
		return nil, err
	}
	// Only get daily notes
	matches, err := fs.Glob(os.DirFS(d), "*-*-*.md")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func execCommand() error {
	cmd := exec.Command("zsh")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, err = in.Write([]byte("open .\n"))
	if err != nil {
		return err
	}
	err = in.Close()
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
