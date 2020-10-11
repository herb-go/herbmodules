package captcha

import (
	"bytes"
	"math/rand"
	"time"
)

type Wanted struct {
	Min          int
	Max          int
	OptionalByte []byte
}

func NewWanted() *Wanted {
	return &Wanted{}
}
func (w *Wanted) DoCaptcha(c Context) (bool, error) {
	data, err := c.GetCaptchaData(ContextNameWanted)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, nil
	}
	submited, err := c.GetCaptchaData(ContextNameSubmited)
	if err != nil {
		return false, err
	}
	if submited == nil {
		return false, nil
	}
	return data.Equal(submited), nil
}

func (w *Wanted) Challenge(c Context) error {
	data, err := c.GetCaptchaData(ContextNameWanted)
	if err != nil {
		return err
	}
	if data != nil {
		return nil
	}
	return c.SetCaptchaData(ContextNameWanted, w.NewWantedBytes())
}

func (w *Wanted) NewWantedBytes() []byte {
	length := w.Min
	if w.Min != w.Max {
		length = length + rand.Intn(w.Max-w.Min)
	}
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = w.OptionalByte[rand.Intn(len(w.OptionalByte))]
	}
	return result
}

func (w *Wanted) ToIndexBytes(data []byte) ([]byte, error) {
	var result = make([]byte, len(data))
	for k := range data {
		index := bytes.IndexByte(w.OptionalByte, data[k])
		if index < 0 {
			return nil, ErrByteNotAvaliable
		}
		result[k] = byte(index)
	}
	return result, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
