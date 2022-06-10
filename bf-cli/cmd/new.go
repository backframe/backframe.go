package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ndaba1/gommander"
)

type Preset struct {
	Interfaces   []string `json:"interfaces"`
	Internals    []string `json:"internals"`
	Integrations []string `json:"integrations"`
	Providers    []string `json:"providers"`
	Database     string   `json:"database"`
}

func DefaultPreset() Preset {
	return Preset{
		Interfaces: []string{
			"Rest",
		},
		Providers: []string{
			"EmailAndPassword",
		},
		Database:     "MongoDB",
		Internals:    []string{},
		Integrations: []string{},
	}
}

func New(pm *gommander.ParserMatches) {
	name, _ := pm.GetArgValue("<app-name>")
	git := pm.ContainsFlag("--git")
	name, targetDir := getProjectName(name, pm)

	// preset to use
	var config Preset

	// set default if prompts skipped
	if pm.ContainsFlag("--default") {
		config = DefaultPreset()
	} else if presetPath, err := pm.GetOptionValue("--preset"); err == nil {
		_, err := os.Stat(presetPath)
		if err != nil {
			panic("An error occurred while trying to load the preset")
		}
		contents, _ := os.ReadFile(presetPath)
		err = json.Unmarshal(contents, &config)
		if err != nil {
			panic("An error occurred while trying to load the preset")
		}
	} else {
		// resolve prompts
		config = resolvePrompts()
	}

	// TODO: Check for save preset, save preset with name in configured location
	file, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(filepath.Join(targetDir, "bfpreset.json"), file, fs.FileMode(os.O_RDWR))

	// Start actual tasks execution
	initializeProject(name, targetDir, config)

	if git {
		initializeGit(targetDir)
	}

	// TODO: print `done` messages
}

func getProjectName(name string, pm *gommander.ParserMatches) (string, string) {
	inCurrent := name == "."

	dir, err := os.Getwd()
	if err != nil {
		panic("An error occurred while getting the current working directory")
	}

	targetDir := dir
	if inCurrent {
		name = filepath.Base(dir)
	} else {
		targetDir = path.Join(dir, name)
	}

	// if dir exists and force flag is not set
	if _, err := os.Stat(targetDir); err == nil {
		if pm.ContainsFlag("--force") {
			err := os.RemoveAll(targetDir)
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
					os.Exit(1)
				}
			} else {
				action := ""
				prompt := &survey.Select{
					Message: "Target directory already exists choose an action",
					Options: []string{"Overwrite existing directory", "Enter a different name", "Cancel"},
				}
				survey.AskOne(prompt, &action)

				if action == "Cancel" || len(action) == 0 {
					os.Exit(1)
				} else if action == "Enter a different name" {
					new := ""
					p := &survey.Input{
						Message: "Enter the new name for your project:",
					}
					survey.AskOne(p, &new)
					return getProjectName(new, pm)
				} else {
					err := os.RemoveAll(targetDir)
					if err != nil {
						fmt.Println(err.Error())
						panic("An error occurred while removing the target directory")
					}
				}
			}
		}
	}

	return name, targetDir
}

func resolvePrompts() Preset {
	return Preset{}
}

func initializeProject(name string, dir string, cfg Preset) {
	// Create project structure
	fmter := gommander.NewFormatter(gommander.DefaultTheme())

	fmter.Add(gommander.Headline, "Generating new project...")
	fmter.Print()

	os.Mkdir(dir, os.ModeAppend)
	// Get project dependencies resolved from the preset
	// Invoke generators accordingly
}

func initializeGit(dir string) {
	// try init git
}
