package jobrole

import "testing"

func TestAssignment(t *testing.T) {
	a := NewAssignment()
	if a.ID != "" || a.JobID != "" || a.After != nil || a.Before != nil || len(a.Data) != 0 {
		t.Fail()
	}
	after := int64(1)
	before := int64(2)
	a.MergeID("assignmentid").MergeJobID("jobid").MergeAfter(&after).MergeBefore(&before).MergeData("name", "value")
	if a.ID != "assignmentid" || a.JobID != "jobid" || a.After != &after || a.Before != &before || len(a.Data) != 1 || a.Data["name"] != "value" {
		t.Fail()
	}
}

func TestNopAssignmentService(t *testing.T) {
	nas := NopAssignmentService{}
	if a, err := nas.GetUserAssignments(""); a != nil || err != nil {
		t.Fail()
	}
}

func TestAsiignmentMap(t *testing.T) {
	am := NewAssignmentMap()
	assignment1 := NewAssignment()
	assignment2 := NewAssignment()
	am.Assign("userid", []*Assignment{assignment1, assignment2})
	if a, err := am.GetUserAssignments(""); len(a) != 0 || err != nil {
		t.Fail()
	}
	if a, err := am.GetUserAssignments("userid"); len(a) != 2 || a[0] != assignment1 || a[1] != assignment2 || err != nil {
		t.Fail()
	}
	am.Assign("userid", nil)
	if a, err := am.GetUserAssignments("userid"); len(a) != 0 || err != nil {
		t.Fail()
	}
	am.Assign("user2", []*Assignment{assignment2})
	am.Assign("user1", []*Assignment{assignment1})
	exported := am.Export()
	if len(exported) != 2 || exported[0].UID != "user1" || len(exported[0].Assignments) != 1 || exported[0].Assignments[0] != assignment1 || exported[1].UID != "user2" || len(exported[1].Assignments) != 1 || exported[1].Assignments[0] != assignment2 {
		t.Fail()
	}
	am.Assign("user2", nil)
	am.Assign("user1", nil)
	if a, err := am.GetUserAssignments("user2"); len(a) != 0 || err != nil {
		t.Fail()
	}
	if a, err := am.GetUserAssignments("user1"); len(a) != 0 || err != nil {
		t.Fail()
	}
	am.Import(exported)
	if a, err := am.GetUserAssignments("user1"); len(a) != 1 || a[0] != assignment1 || err != nil {
		t.Fail()
	}
	if a, err := am.GetUserAssignments("user2"); len(a) != 1 || a[0] != assignment2 || err != nil {
		t.Fail()
	}
}
