package jobrole

import (
	"sort"
	"sync"
)

type Assignment struct {
	ID       string
	UserID   string
	JobID    string
	CacheKey string
	After    *int64
	Before   *int64
	Data     map[string]string
}

func NewAssignment() *Assignment {
	return &Assignment{}
}

type Assignments []*Assignment

func (a Assignments) Len() int {
	return len(a)
}
func (a Assignments) Less(i, j int) bool {
	return a[i].ID < a[j].ID
}
func (a Assignments) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type AssignmentService interface {
	Assign(uid string, assignment *Assignment) error
	RemoveAssignment(assignmentid string) error
	GetUserAssignments(uid string) ([]*Assignment, error)
}
type NopAssignmentService struct{}

func (s NopAssignmentService) Assign(uid string, assignment *Assignment) error {
	return nil
}
func (s NopAssignmentService) RemoveAssignment(assignmentid string) error {
	return nil
}

func (s NopAssignmentService) GetUserAssignments(uid string) ([]*Assignment, error) {
	return nil, nil
}

type AssignmentMap struct {
	NopAssignmentService
	byUser map[string][]*Assignment
	byID   map[string]*Assignment
	locker sync.Mutex
}

func NewAssignmentMap() *AssignmentMap {
	return &AssignmentMap{}
}
func (m *AssignmentMap) assign(a *Assignment) {
	if a != nil {
		m.remove(a.ID)
		m.byID[a.ID] = a
		userAssigments := m.byUser[a.UserID]
		if userAssigments == nil {
			userAssigments = []*Assignment{}
		}
		userAssigments = append(userAssigments, a)
		m.byUser[a.UserID] = userAssigments
	}
}
func (m *AssignmentMap) remove(assignmentID string) {
	old := m.byID[assignmentID]
	if old != nil {
		userAssignments := m.byUser[old.UserID]
		if userAssignments != nil {
			result := []*Assignment{}
			for _, v := range userAssignments {
				if v.ID != assignmentID {
					result = append(result, v)
				}
			}
			if len(result) > 0 {
				m.byUser[old.UserID] = result
			} else {
				delete(m.byUser, old.UserID)
			}
		}
		delete(m.byID, assignmentID)
	}
}
func (m *AssignmentMap) Assign(assignment *Assignment) error {
	defer m.locker.Unlock()
	m.locker.Lock()
	m.remove(assignment.ID)
	m.assign(assignment)
	return nil
}

func (m *AssignmentMap) RemoveAssignment(assignmentid string) error {
	defer m.locker.Unlock()
	m.locker.Lock()
	m.remove(assignmentid)
	return nil
}
func (m *AssignmentMap) GetUserAssignments(uid string) ([]*Assignment, error) {
	defer m.locker.Unlock()
	m.locker.Lock()
	ua := m.byUser[uid]
	return ua, nil
}

func (m *AssignmentMap) List() []*Assignment {
	defer m.locker.Unlock()
	m.locker.Lock()
	result := make([]*Assignment, 0, len(m.byID))
	for _, v := range m.byID {
		result = append(result, v)
	}
	sort.Sort(Assignments(result))
	return result
}
