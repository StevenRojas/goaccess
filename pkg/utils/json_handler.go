package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

// JSONHandler interface
type JSONHandler interface {
	Modules() ([]entities.ModuleInit, error)
}

type jsonHandler struct {
	folder string
}

// NewJSONHandler json handler instance
func NewJSONHandler(folderPath string) JSONHandler {
	return &jsonHandler{
		folder: folderPath,
	}
}

// Modules read json files from init folder and return a list of modules
func (jh *jsonHandler) Modules() ([]entities.ModuleInit, error) {
	var files []string
	err := filepath.Walk(jh.folder+"/modules", collect(&files))
	if err != nil {
		return nil, err
	}
	modules := []entities.ModuleInit{}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		module := entities.ModuleInit{}
		err = json.Unmarshal([]byte(content), &module)
		if err != nil {
			return nil, err
		}
		for si, submodule := range module.SubModules {
			sections := make(map[string]bool)
			for _, section := range submodule.SectionList {
				sections[section] = false
			}
			submodule.Sections = sections
			submodule.SectionList = nil

			actions := make(map[string]entities.Action)
			for a, title := range submodule.ActionList {
				actions[a] = entities.Action{
					Title:   title,
					Allowed: false,
				}
			}
			submodule.Actions = actions
			submodule.ActionList = nil

			module.SubModules[si] = submodule
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func collect(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			*files = append(*files, path)
		}
		return nil
	}
}
