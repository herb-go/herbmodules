package jobrole

import (
	"testing"
	"time"
)

func TestService(t *testing.T) {
	am := NewAssignmentMap()
	jm := NewJobMap()
	dm := NewDutyMap()
	s := New()
	s.AssignmentService = am
	s.JobService = jm
	s.DutyService = dm
	t1 := int64(1)
	t2 := int64(3)
	dm.SetDuty(NewDuty().MergeID("duty1").MergeRoles("role1"))
	dm.SetDuty(NewDuty().MergeID("duty2").MergeRoles("role2"))
	dm.SetDuty(NewDuty().MergeID("duty3").MergeRoles("role3"))
	dm.SetDuty(NewDuty().MergeID("duty4").MergeRoles("role4"))
	jm.SetJob(NewJob().MergeID("job1").AppendDuty("duty1"))
	jm.SetJob(NewJob().MergeID("job2").AppendDuty("duty2"))
	jm.SetJob(NewJob().MergeID("job3").AppendDuty("duty3"))
	jm.SetJob(NewJob().MergeID("job4").AppendDuty("duty4"))
	u1a := []*Assignment{
		NewAssignment().MergeJobID("job1"),
		NewAssignment().MergeJobID("job2").MergeAfter(&t1),
		NewAssignment().MergeJobID("job3").MergeBefore(&t2),
		NewAssignment().MergeJobID("job3").MergeBefore(&t1).MergeAfter(&t2),
	}
	am.Assign("user1", u1a)
	time1 := int64(2)
	r, err := s.GetUserRoles("user1", &time1)
	data := r.Data()
	if len(data) != 3 || err != nil {
		t.Fatal()
	}
	time2 := int64(0)
	r, err = s.GetUserRoles("user1", &time2)
	data = r.Data()
	if len(data) != 2 || err != nil {
		t.Fatal()
	}
	time3 := int64(5)
	r, err = s.GetUserRoles("user1", &time3)
	data = r.Data()
	if len(data) != 2 || err != nil {
		t.Fatal()
	}
	r, err = s.GetUserRoles("notexist", &time3)
	data = r.Data()
	if len(data) != 0 || err != nil {
		t.Fatal()
	}
	t3 := time.Now().Add(-10 * time.Minute).Unix()
	t4 := time.Now().Add(10 * time.Minute).Unix()
	u2a := []*Assignment{
		NewAssignment().MergeJobID("job1"),
		NewAssignment().MergeJobID("job2").MergeAfter(&t3),
		NewAssignment().MergeJobID("job3").MergeBefore(&t4),
		NewAssignment().MergeJobID("job3").MergeBefore(&t3).MergeAfter(&t4),
	}
	am.Assign("user1", u2a)
	r, err = s.GetUserRoles("user1", nil)
	data = r.Data()
	if len(data) != 3 || err != nil {
		t.Fatal()
	}
}
