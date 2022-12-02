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

// New return Template instance
func New(templateName string, forwardBody []byte) (*Template, error) {
	templatesPath := getTemplatePath()
	if err := checkTemplateName(templatesPath, templateName); err != nil {
		return nil, err
	}
	fileBody, err := os.ReadFile(filepath.Join(templatesPath, templateName))
	if err != nil {
		return nil, err
	}
	return &Template{
		body:        string(fileBody),
		forwardBody: forwardBody,
	}, nil
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

func getTemplatePath() string {
	userHomeDir := berrors.Unwrap(os.UserHomeDir())
	templatesPath := filepath.Join(userHomeDir, "/.config/webhook-forwarder/templates")
	return templatesPath
}

// JsResult created after Template.RunJs
type JsResult struct {
	Target  string
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
		convertValue = berrors.Unwrap(vm.Call("convert", nil, t.forwardBody))
	} else {
		convertValue = berrors.Unwrap(vm.Call("convert", nil, r))
	}
	if !convertValue.IsObject() {
		return nil, fmt.Errorf("js template incorrent")
	}
	target, err := getTarget(&convertValue)
	if err != nil {
		return nil, err
	}
	payload, err := GetObject(&convertValue)
	if err != nil {
		return nil, err
	}
	return &JsResult{target, payload}, nil
}

func getTarget(v *otto.Value) (string, error) {
	v2, err := v.Object().Get("target")
	if err != nil {
		return "", fmt.Errorf("get value from otto object failed: %+v", err)
	}
	return v2.String(), nil
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
