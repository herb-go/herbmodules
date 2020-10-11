package auth

import "testing"

func TestProfile(t *testing.T) {
	p := &Profile{}
	p.SetValue(ProfileIndexID, "test")
	re := p.Value(ProfileIndexID)
	if re != "test" {
		t.Error(re)
	}
	p.AddValue(ProfileIndexID, "test1")
	re = p.Value(ProfileIndexID)
	if re != "test" {
		t.Error(re)
	}
	res := p.Values(ProfileIndexID)
	if len(res) != 2 || res[0] != "test" || res[1] != "test1" {
		t.Error(res)
	}
	p.SetValues(ProfileIndexID, []string{"test3", "test4"})
	re = p.Value(ProfileIndexID)
	if re != "test3" {
		t.Error(re)
	}
	res = p.Values(ProfileIndexID)
	if len(res) != 2 || res[0] != "test3" || res[1] != "test4" {
		t.Error(res)
	}
	p.AddValue(ProfileIndexEmail, "testadd")
	re = p.Value(ProfileIndexEmail)
	if re != "testadd" {
		t.Error(re)
	}
	re = p.Value(ProfileIndexFirstName)
	if re != "" {
		t.Error(re)
	}
	res = p.Values(ProfileIndexFirstName)
	if res != nil {
		t.Error(res)
	}
}
