package captcha

import "testing"

func TestCaptcha(t *testing.T) {
	var err error
	var result bool
	c := New()
	c.Driver = NewWanted()
	ctx := NewPlainContext()
	err = ctx.SetCaptchaData(ContextNameWanted, []byte("Value1"))
	if err != nil {
		t.Fatal(err)
	}
	err = ctx.SetCaptchaData(ContextNameSubmited, []byte("Value2"))
	if err != nil {
		t.Fatal(err)
	}
	result, err = c.DoCaptcha(ctx)
	if result != false || err != nil {
		t.Fatal(result, err)
	}
	ctx.Trusted = true
	result, err = c.DoCaptcha(ctx)
	if result != true || err != nil {
		t.Fatal(result, err)
	}
}
