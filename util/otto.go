package util

import (
	"encoding/json"
	"fmt"
	"github.com/Bpazy/berrors"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
)

type Template struct {
	body        string
	forwardBody []byte
}

func NewTemplate(body string, forwardBody []byte) *Template {
	return &Template{
		body:        body,
		forwardBody: forwardBody,
	}
}

type RunResult struct {
	otto.Value
}

func (t Template) RunJs() (*RunResult, error) {
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
	return &RunResult{Value: convertValue}, nil
}

// GetString the value string of the property with the given name.
func (t RunResult) GetString(name string) (string, error) {
	v, err := t.Object().Get(name)
	if err != nil {
		return "", fmt.Errorf("get value from otto object failed: %+v", err)
	}
	return v.String(), nil
}

// GetObject the value object of the property with the given name.
func (t RunResult) GetObject(name string) (any, error) {
	v, err := t.Object().Get(name)
	if err != nil {
		return nil, fmt.Errorf("get value from otto object failed: %+v", err)
	}
	r, err := v.Export()
	if err != nil {
		return nil, fmt.Errorf("get object from otto value failed: %+v", err)
	}
	return r, nil
}
