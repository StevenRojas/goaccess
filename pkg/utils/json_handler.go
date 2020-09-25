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
	Modules() ([]entities.Module, error)
	Actions() ([]entities.ActionModule, error)
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
func (jh *jsonHandler) Modules() ([]entities.Module, error) {
	var files []string
	err := filepath.Walk(jh.folder+"/modules", collect(&files))
	if err != nil {
		return nil, err
	}
	modules := []entities.Module{}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		module := entities.Module{}
		err = json.Unmarshal([]byte(content), &module)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func (jh *jsonHandler) Actions() ([]entities.ActionModule, error) {
	var files []string
	err := filepath.Walk(jh.folder+"/actions", collect(&files))
	if err != nil {
		return nil, err
	}
	actions := []entities.ActionModule{}
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		action := entities.ActionModule{}
		err = json.Unmarshal([]byte(content), &action)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
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
