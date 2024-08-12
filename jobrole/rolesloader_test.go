package jobrole

import (
	"testing"

	"github.com/herb-go/herbsecurity/authorize/role/roleparser"
)

func TestRolesLoader(t *testing.T) {
	s := New()
	dm := NewDutyMap()
	jm := NewJobMap()
	s.DutyService = dm
	s.JobService = jm
	duty := NewDuty()
	duty.ID = "duty"
	duty.Roles = "test"
	duty2 := NewDuty()
	duty2.ID = "duty2"
	duty2.Roles = "testreplacer"
	dutynotexist := NewDuty()
	dutynotexist.ID = "dutynotexist"
	dutynotexist.Roles = "testreplacer"

	dm.SetDuty(duty)
	dm.SetDuty(duty2)
	job := NewJob().MergeID("job").AppendDuty(duty.ID)
	job2 := NewJob().MergeID("job2").AppendDuty(duty2.ID)
	jobdutynotexsit := NewJob().MergeID("jobdutynotexsit").AppendDuty(dutynotexist.ID)
	jm.SetJob(job)
	jm.SetJob(job2)
	jm.SetJob(jobdutynotexsit)
	assignment := NewAssignment()
	assignment.JobID = job.ID
	assignment.Data["testreplacer"] = "testvalue"
	assignment2 := NewAssignment()
	assignment2.JobID = job2.ID
	assignment2.Data["testreplacer"] = "testvalue"
	assignmentdutynotexsit := NewAssignment()
	assignmentdutynotexsit.JobID = jobdutynotexsit.ID
	assignmentnotexsit := NewAssignment()
	assignmentnotexsit.JobID = "notexist"

	replacedRoles := "testvalue"
	pl := PlainRolesLoader{}
	if r, err := pl.LoadRoles(s, nil); r != nil || err != nil {
		t.Fatal()
	}
	if r, err := pl.LoadRoles(s, assignmentdutynotexsit); len(r) != 0 || err != nil {
		t.Fatal()
	}
	if r, err := pl.LoadRoles(s, assignmentnotexsit); len(r) != 0 || err != nil {
		t.Fatal()
	}
	if r, err := pl.LoadRoles(s, assignment); err != nil || len(r) == 0 || roleparser.StringifyRoles(r[0]) != duty.Roles {
		t.Fatal()
	}
	rl := ReplacerRolesLoader{}
	if r, err := rl.LoadRoles(s, nil); r != nil || err != nil {
		t.Fatal()
	}
	if r, err := rl.LoadRoles(s, assignmentdutynotexsit); len(r) != 0 || err != nil {
		t.Fatal()
	}
	if r, err := rl.LoadRoles(s, assignmentnotexsit); len(r) != 0 || err != nil {
		t.Fatal()
	}

	if r, err := rl.LoadRoles(s, assignment); err != nil || len(r) == 0 || roleparser.StringifyRoles(r[0]) != duty.Roles {
		t.Fatal()
	}
	if r, err := rl.LoadRoles(s, assignment2); err != nil || len(r) == 0 || roleparser.StringifyRoles(r[0]) == duty2.Roles || roleparser.StringifyRoles(r[0]) != replacedRoles {
		t.Fatal()
	}
}
