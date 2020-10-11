package captcha

import "testing"

func TestWanted(t *testing.T) {
	var c Context
	var result bool
	var err error
	var wrongData = []byte("12345")
	var wantedData ContextData
	var data ContextData
	w := NewWanted()
	w.Min = 8
	w.Max = 10
	w.OptionalByte = []byte("abcde")
	c = NewPlainContext()
	result, err = w.DoCaptcha(c)
	if result != false || err != nil {
		t.Fatal(w)
	}
	err = c.SetCaptchaData(ContextNameWanted, wrongData)
	if err != nil {
		t.Fatal(w)
	}
	result, err = w.DoCaptcha(c)
	if result != false || err != nil {
		t.Fatal(w)
	}
	c = NewPlainContext()
	err = w.Challenge(c)
	if err != nil {
		t.Fatal(w)
	}
	wantedData, err = c.GetCaptchaData(ContextNameWanted)
	if wantedData == nil || err != nil {
		t.Fatal(c)
	}
	err = w.Challenge(c)
	if err != nil {
		t.Fatal(w)
	}
	data, err = c.GetCaptchaData(ContextNameWanted)
	if !data.Equal(wantedData) || err != nil {
		t.Fatal(c)
	}
	c.SetCaptchaData(ContextNameSubmited, wantedData)
	result, err = w.DoCaptcha(c)
	if result != true || err != nil {
		t.Fatal(w)
	}
}
