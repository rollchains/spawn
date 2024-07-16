package main

import (
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/charmbracelet/glow/ui"
	"github.com/rollchains/spawn/docs"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var DocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "View Spawn Documentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath := path.Join(os.TempDir(), "spawn", "docs")

		err := fs.WalkDir(docs.Docs, ".", func(relPath string, d fs.DirEntry, e error) error {
			newPath := path.Join(dirPath, relPath)

			if d.IsDir() {
				return os.MkdirAll(newPath, 0755)
			}

			fi, err := docs.Docs.Open(relPath)
			if err != nil {
				return err
			}
			defer fi.Close()

			fo, err := os.Create(newPath)
			if err != nil {
				return err
			}
			defer fo.Close()

			_, err = io.Copy(fo, fi)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
		defer os.RemoveAll(dirPath)

		return View(dirPath)
	},
}

// orginal source: https://github.com/ignite/cli/blob/main/ignite/pkg/markdownviewer/markdownviewer.go
func View(path string) error {
	conf, err := config(path)
	if err != nil {
		return err
	}

	_, err = ui.NewProgram(conf).Run()
	return err
}

func config(path string) (ui.Config, error) {
	var width uint

	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return ui.Config{}, err
	}
	width = uint(w)
	if width > 120 {
		width = 120
	}

	docTypes := ui.NewDocTypeSet()
	docTypes.Add(ui.LocalDoc)

	conf := ui.Config{
		WorkingDirectory:     path,
		DocumentTypes:        docTypes,
		GlamourStyle:         "auto",
		HighPerformancePager: true,
		GlamourEnabled:       true,
		GlamourMaxWidth:      width,
	}

	return conf, nil
}
