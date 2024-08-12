package jobrole

import (
	"sort"
	"sync"
)

type Assignment struct {
	ID     string
	JobID  string
	After  *int64
	Before *int64
	Data   map[string]string
}

func (a *Assignment) MergeID(id string) *Assignment {
	a.ID = id
	return a
}
func (a *Assignment) MergeJobID(jobid string) *Assignment {
	a.JobID = jobid
	return a
}
func (a *Assignment) MergeAfter(after *int64) *Assignment {
	a.After = after
	return a
}
func (a *Assignment) MergeBefore(before *int64) *Assignment {
	a.Before = before
	return a
}
func (a *Assignment) MergeData(name string, value string) *Assignment {
	a.Data[name] = value
	return a
}
func NewAssignment() *Assignment {
	return &Assignment{
		Data: map[string]string{},
	}
}

type AssignmentService interface {
	GetUserAssignments(uid string) ([]*Assignment, error)
}
type NopAssignmentService struct{}

func (s NopAssignmentService) GetUserAssignments(uid string) ([]*Assignment, error) {
	return nil, nil
}

type AssignmentMap struct {
	NopAssignmentService
	data   map[string][]*Assignment
	locker sync.Mutex
}

func NewAssignmentMap() *AssignmentMap {
	return &AssignmentMap{
		data: map[string][]*Assignment{},
	}
}
func (m *AssignmentMap) Assign(uid string, assignments []*Assignment) {
	defer m.locker.Unlock()
	m.locker.Lock()
	if len(assignments) == 0 {
		delete(m.data, uid)
	} else {
		m.data[uid] = assignments
	}
}

func (m *AssignmentMap) GetUserAssignments(uid string) ([]*Assignment, error) {
	defer m.locker.Unlock()
	m.locker.Lock()
	ua := m.data[uid]
	if ua == nil {
		ua = []*Assignment{}
	}
	return ua, nil
}

type UserAssignment struct {
	UID         string
	Assignments []*Assignment
}

func (m *AssignmentMap) Export() []*UserAssignment {
	defer m.locker.Unlock()
	m.locker.Lock()
	result := make([]*UserAssignment, 0, len(m.data))
	for k, v := range m.data {
		result = append(result, &UserAssignment{UID: k, Assignments: append([]*Assignment{}, v...)})
	}
	sort.Sort(UserAssignments(result))
	return result
}

func (m *AssignmentMap) Import(data []*UserAssignment) {
	defer m.locker.Unlock()
	m.locker.Lock()
	for _, v := range data {
		m.data[v.UID] = append([]*Assignment{}, v.Assignments...)
	}
}

type UserAssignments []*UserAssignment

func (a UserAssignments) Len() int {
	return len(a)
}

func (a UserAssignments) Less(i, j int) bool {
	return a[i].UID < a[j].UID
}

func (a UserAssignments) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
