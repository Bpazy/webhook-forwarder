package template

import (
	"encoding/json"
	"fmt"
	"github.com/Bpazy/berrors"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"os"
	"path/filepath"
)

type Template struct {
	body        string
	forwardBody []byte
}

func GetTemplateContent(templateName string) (string, error) {
	templatesPath := getTemplatePath()
	if err := checkTemplateName(templatesPath, templateName); err != nil {
		return "", err
	}
	content, err := os.ReadFile(filepath.Join(templatesPath, templateName))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// New return Template instance
func New(js string, forwardBody []byte) (*Template, error) {
	return &Template{
		body:        js,
		forwardBody: forwardBody,
	}, nil
}

func getTemplatePath() string {
	userHomeDir := berrors.Unwrap(os.UserHomeDir())
	templatesPath := filepath.Join(userHomeDir, "/.config/webhook-forwarder/templates")
	return templatesPath
}

func checkTemplateName(templatesPath string, name string) error {
	dirs := berrors.Unwrap(os.ReadDir(templatesPath))
	for _, dir := range dirs {
		if dir.Name() == name {
			return nil
		}
	}
	return fmt.Errorf("no template named %s", name)
}

// JsResult created after Template.RunJs
type JsResult struct {
	Targets []string
	Payload any
}

// RunJs run js from template and return result
func (t Template) RunJs() (*JsResult, error) {
	vm := otto.New()
	if _, err := vm.Run(t.body); err != nil {
		return nil, err
	}

	r := gin.H{}
	var convertValue otto.Value
	if err := json.Unmarshal(t.forwardBody, &r); err != nil {
		convertValue, err = vm.Call("convert", nil, t.forwardBody)
		if err != nil {
			return nil, fmt.Errorf("call js function 'convert' error: %+v", err)
		}
	} else {
		convertValue, err = vm.Call("convert", nil, r)
		if err != nil {
			return nil, fmt.Errorf("call js function 'convert' error: %+v", err)
		}
	}
	if !convertValue.IsObject() {
		return nil, fmt.Errorf("js template incorrent")
	}
	targets, err := getTargets(&convertValue)
	if err != nil {
		return nil, err
	}
	payload, err := GetObject(&convertValue)
	if err != nil {
		return nil, err
	}
	return &JsResult{targets, payload}, nil
}

var UnsupportedTargetType = fmt.Errorf("unsupported target type")

func getTargets(v *otto.Value) ([]string, error) {
	v2, err := v.Object().Get("target")
	if err != nil {
		return []string{}, fmt.Errorf("get target from otto object failed: %+v", err)
	}
	i, err := v2.Export()
	if err != nil {
		return []string{}, fmt.Errorf("get target from otto object failed: %+v", err)
	}
	switch i.(type) {
	case string:
		return []string{i.(string)}, nil
	case []string:
		return i.([]string), nil
	}
	return []string{}, UnsupportedTargetType
}

func GetObject(v *otto.Value) (any, error) {
	v2, err := v.Object().Get("payload")
	if err != nil {
		return nil, fmt.Errorf("get value from otto object failed: %+v", err)
	}
	r, err := v2.Export()
	if err != nil {
		return nil, fmt.Errorf("get object from otto value failed: %+v", err)
	}
	return r, nil
}
