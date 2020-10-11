package captcha

import (
	"bytes"
	"testing"
)

func TestContext(t *testing.T) {
	var data ContextData
	var err error
	var scene Scene
	var trusted bool
	var SceneTest = Scene("testscene")
	var TestData = []byte("testdata")
	var c = NewPlainContext()
	data, err = c.GetCaptchaData(ContextNameSubmited)
	if data != nil || err != nil {
		t.Fatal(c)
	}
	err = c.SetCaptchaData(ContextNameSubmited, TestData)
	if err != nil {
		t.Fatal(c)
	}
	data, err = c.GetCaptchaData(ContextNameSubmited)
	if bytes.Compare(data, TestData) != 0 || err != nil {
		t.Fatal(c)
	}
	scene, err = c.CaptchaScene()
	if scene != DefaultScene || err != nil {
		t.Fatal(c)
	}
	c.Scene = SceneTest
	scene, err = c.CaptchaScene()
	if scene != SceneTest || err != nil {
		t.Fatal(c)
	}

	trusted, err = c.CaptchaTrusted()
	if trusted != false || err != nil {
		t.Fatal(c)
	}
	c.Trusted = true
	trusted, err = c.CaptchaTrusted()
	if trusted != true || err != nil {
		t.Fatal(c)
	}
}
