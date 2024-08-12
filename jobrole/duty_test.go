package jobrole

import "testing"

func TestDuty(t *testing.T) {
	d := NewDuty()
	if d.ID != "" || d.Title != "" || d.Desc != "" || d.Roles != "" {
		t.Fail()
	}
	d.MergeID("testid").MergeTitle("testtitle").MergeDesc("testdesc").MergeRoles("testroles")
	if d.ID != "testid" || d.Title != "testtitle" || d.Desc != "testdesc" || d.Roles != "testroles" {
		t.Fail()
	}
}
func TestNopDutyService(t *testing.T) {
	nds := NopDutyService{}
	if job, err := nds.GetDuty("test"); job != nil || err != nil {
		t.Fail()
	}
}

func TestDutyMap(t *testing.T) {
	dm := NewDutyMap()
	duty1 := NewDuty().MergeID("duty1")
	duty2 := NewDuty().MergeID("duty2")
	if list := dm.List(); len(list) != 0 {
		t.Fail()
	}
	dm.SetDuty(duty2)
	dm.SetDuty(duty1)
	if d, err := dm.GetDuty("notexist"); d != nil || err != nil {
		t.Fail()
	}
	if d, err := dm.GetDuty(duty1.ID); d != duty1 || err != nil {
		t.Fail()
	}
	if d, err := dm.GetDuty(duty2.ID); d != duty2 || err != nil {
		t.Fail()
	}
	if list := dm.List(); len(list) != 2 || list[0] != duty1 || list[1] != duty2 {
		t.Fatal()
	}
}
