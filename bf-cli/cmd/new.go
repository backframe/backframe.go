package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ndaba1/gommander"
)

func New(pm *gommander.ParserMatches) {
	name, _ := pm.GetArgValue("<app-name>")
	// git := pm.ContainsFlag("--git")
	inCurrent := name == "."

	dir, err := os.Getwd()
	if err != nil {
		panic("An error occurred while getting the current working directory")
	}

	target_dir := dir
	if inCurrent {
		name = filepath.Base(dir)
	} else {
		target_dir = path.Join(dir, name)
	}

	// if dir exists and force flag is not set
	if _, err := os.Stat(target_dir); err == nil {
		if pm.ContainsFlag("--force") {
			err := os.RemoveAll(target_dir)
			if err != nil {
				panic("An error occurred while removing the target directory")
			}
		} else {
			if inCurrent {
				overwrite := false
				prompt := &survey.Confirm{
					Message: "Generate project in current directory?",
				}
				survey.AskOne(prompt, &overwrite)

				if !overwrite {
					return
				}
			} else {
				action := ""
				prompt := &survey.Select{
					Message: "Target directory already exists choose an action",
					Options: []string{"Overwrite existing directory", "Cancel"},
				}
				survey.AskOne(prompt, &action)

				if action == "Cancel" || len(action) == 0 {
					return
				} else {
					err := os.RemoveAll(target_dir)
					if err != nil {
						fmt.Println(err.Error())
						panic("An error occurred while removing the target directory")
					}
				}
			}
		}
	}
}
