package jobrole

import "sync"

type Duty struct {
	ID    string
	Title string
	Desc  string
	Roles string
}

func NewDuty() *Duty {
	return &Duty{}
}

type DutyMap struct {
	NopDutyServive
	data map[string]*Duty
	lock sync.Mutex
}

func NewDutyMap() *DutyMap {
	return &DutyMap{}
}

func (m *DutyMap) GetDuty(id string) (*Duty, error) {
	defer m.lock.Unlock()
	m.lock.Lock()
	return m.data[id], nil
}
func (m *DutyMap) SetDuty(id string, duty *Duty) error {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.data[id] = duty
	return nil
}
func (m *DutyMap) List() []*Duty {
	defer m.lock.Unlock()
	m.lock.Lock()

	result := make([]*Duty, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}
	return result
}

type DutyService interface {
	GetDuty(id string) (*Duty, error)
	SetDuty(id string, duty *Duty) error
}

type NopDutyServive struct{}

func (s NopDutyServive) GetDuty(id string) (*Duty, error) {
	return nil, nil
}

func (s NopDutyServive) SetDuty(id string, duty *Duty) error {
	return nil
}
