package jobrole

import (
	"strings"
	"testing"
)

func TestJob(t *testing.T) {
	j := NewJob()
	if j.ID != "" || j.Title != "" || j.Desc != "" || len(j.DutyList) != 0 {
		t.Fail()
	}
	j.MergeID("testid").MergeTitle("testtitle").MergeDesc("testdesc").AppendDuty("duty1").AppendDuty("duty2")
	if j.ID != "testid" || j.Title != "testtitle" || j.Desc != "testdesc" || len(j.DutyList) != 2 || strings.Join(j.DutyList, ",") != "duty1,duty2" {
		t.Fail()
	}
}
func TestNopJobService(t *testing.T) {
	njs := NopJobService{}
	if job, err := njs.GetJob("test"); job != nil || err != nil {
		t.Fail()
	}
}

func TestJobMap(t *testing.T) {
	jm := NewJobMap()
	job := NewJob().MergeID("job1")
	job2 := NewJob().MergeID("job2")
	if list := jm.List(); len(list) != 0 {
		t.Fail()
	}
	jm.SetJob(job2)
	jm.SetJob(job)
	if j, err := jm.GetJob("notexist"); j != nil || err != nil {
		t.Fail()
	}
	if j, err := jm.GetJob(job.ID); j != job || err != nil {
		t.Fail()
	}
	if j, err := jm.GetJob(job2.ID); j != job2 || err != nil {
		t.Fail()
	}
	if list := jm.List(); len(list) != 2 || list[0] != job || list[1] != job2 {
		t.Fail()
	}
}
