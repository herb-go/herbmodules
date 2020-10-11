package captcha

import (
	"bytes"
)

type ContextValue struct {
	Name ContextName
	Data ContextData
}

type ContextData []byte

func (d ContextData) Equal(data ContextData) bool {
	return bytes.Compare(d, data) == 0
}

type Collection map[ContextName]ContextData

func (c *Collection) SetCaptchaData(n ContextName, v ContextData) error {
	(*c)[n] = v
	return nil
}
func (c *Collection) GetCaptchaData(n ContextName) (ContextData, error) {
	return (*c)[n], nil
}

func NewCollection() *Collection {
	return &Collection{}
}
