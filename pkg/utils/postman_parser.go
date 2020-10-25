package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"github.com/StevenRojas/goaccess/pkg/entities"
)

const baseURL = "{{baseUrl}}/" // this pattern will be replaced with an empty string

// PostmanParser postman parser interface
type PostmanParser interface {
	Parse(filename string, module string) (string, error)
}

type postman struct {
	folder  string
	useDesc bool
	module  entities.ModuleInit
}

// NewPostmanParser get a postman parser instance
func NewPostmanParser(folderPath string, useDesc bool) PostmanParser {
	return &postman{
		folder:  folderPath,
		useDesc: useDesc,
	}
}

func (s *postman) Parse(filename string, module string) (string, error) {
	content, err := ioutil.ReadFile(s.folder + "/" + filename)
	if err != nil {
		return "", err
	}
	var j map[string]interface{}
	err = json.Unmarshal([]byte(content), &j)
	if err != nil {
		return "", err
	}
	s.module.Name = module
	submoduleList := make(map[string]bool)
	collection := j["item"].([]interface{})
	for _, c := range collection {
		var actionList = map[string]string{}
		for k, items := range c.(map[string]interface{}) {
			if k == "item" {
				submodule := entities.SubModuleInit{}
				for _, item := range items.([]interface{}) {
					endpoint := item.(map[string]interface{})
					request := endpoint["request"].(map[string]interface{})
					if request["method"] == "GET" {
						continue
					}
					rawURL := request["url"].(map[string]interface{})
					name := reflect.ValueOf(rawURL["path"]).Index(0)
					submoduleName := name.Interface().(string)

					if _, ok := submoduleList[submoduleName]; !ok {
						submodule.Name = submoduleName
						submoduleList[submoduleName] = true
					}
					action, description, ok := s.parseRequest(request)
					if ok {
						actionList[action] = description
					}
				}
				submodule.ActionList = actionList
				submodule.SectionList = []string{"sections here"}
				s.module.SubModules = append(s.module.SubModules, submodule)
			}
		}
	}
	m, _ := json.Marshal(s.module)
	return string(m), nil
}

func (s *postman) parseRequest(request map[string]interface{}) (string, string, bool) {
	rawURL := request["url"].(map[string]interface{})
	url := strings.Replace(rawURL["raw"].(string), baseURL, "", 1)
	re := regexp.MustCompile(`\?.*`) // Find URL params to remove them
	url = re.ReplaceAllString(url, "")
	re = regexp.MustCompile(`\/:([^\/]*?)(\/|$)`) // Find URL plaseholders to replace them
	url = re.ReplaceAllString(url, ":[]:")
	url = strings.TrimRight(url, ":")
	url = strings.Replace(url, "/", ":", -1)
	method := strings.ToLower(request["method"].(string))
	action := method + ":" + url
	var description string
	if s.useDesc {
		desc, ok := request["description"]
		if ok {
			description = fmt.Sprintf("%s", desc)
		} else {
			description = s.getDescription(method, url)
		}
	} else {
		description = s.getDescription(method, url)
	}

	return action, description, true
}

func (s *postman) getDescription(method string, url string) string {
	var prefix string
	switch method {
	case "post":
		prefix = "Create"
	case "put":
		prefix = "Update"
	case "patch":
		prefix = "Modify"
	case "delete":
		prefix = "Delete"
	}
	re := regexp.MustCompile(`\W`) // Remove all non word characters
	suffix := re.ReplaceAllString(url, " ")
	re = regexp.MustCompile(`\s+`) // Replace multiple spaces
	suffix = re.ReplaceAllString(suffix, " ")
	return prefix + " " + strings.Trim(suffix, " ")
}
