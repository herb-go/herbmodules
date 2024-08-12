package jobrole

import (
	"sort"
	"sync"
)

type Duty struct {
	ID    string
	Title string
	Desc  string
	Roles string
}

func NewDuty() *Duty {
	return &Duty{}
}

func (d *Duty) MergeID(id string) *Duty {
	d.ID = id
	return d
}
func (d *Duty) MergeTitle(title string) *Duty {
	d.Title = title
	return d
}
func (d *Duty) MergeDesc(desc string) *Duty {
	d.Desc = desc
	return d
}
func (d *Duty) MergeRoles(roles string) *Duty {
	d.Roles = roles
	return d
}

type DutyMap struct {
	NopDutyService
	data map[string]*Duty
	lock sync.Mutex
}

func NewDutyMap() *DutyMap {
	return &DutyMap{
		data: map[string]*Duty{},
	}
}

func (m *DutyMap) GetDuty(id string) (*Duty, error) {
	defer m.lock.Unlock()
	m.lock.Lock()
	return m.data[id], nil
}
func (m *DutyMap) SetDuty(duty *Duty) {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.data[duty.ID] = duty
}
func (m *DutyMap) List() []*Duty {
	defer m.lock.Unlock()
	m.lock.Lock()
	result := make([]*Duty, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}
	sort.Sort(Duties(result))
	return result
}

type DutyService interface {
	GetDuty(id string) (*Duty, error)
}

type NopDutyService struct{}

func (s NopDutyService) GetDuty(id string) (*Duty, error) {
	return nil, nil
}

type Duties []*Duty

func (d Duties) Len() int {
	return len(d)
}
func (d Duties) Less(i, j int) bool {
	return d[i].ID < d[j].ID
}
func (d Duties) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
