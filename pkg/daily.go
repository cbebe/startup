package startup

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
)

func getDailyNotes() ([]string, error) {
	fsys := os.DirFS(path.Join(os.Getenv("HOME"), "notes", "daily"))
	// Only get daily notes
	matches, err := fs.Glob(fsys, "*-*-*.md")
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
