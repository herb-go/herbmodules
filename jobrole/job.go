package jobrole

import "sync"

type Job struct {
	ID       string
	Title    string
	Desc     string
	DutyList []*Duty
}

func NewJob() *Job {
	return &Job{
		DutyList: []*Duty{},
	}
}

type JobMap struct {
	NopJobService
	data map[string]*Job
	lock sync.Mutex
}

func NewJobMap() *JobMap {
	return &JobMap{}
}

func (m *JobMap) GetJob(id string) (*Job, error) {
	defer m.lock.Unlock()
	m.lock.Lock()
	return m.data[id], nil
}

func (m *JobMap) SetJob(id string, job *Job) error {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.data[id] = job
	return nil
}
func (m *JobMap) List() []*Job {
	defer m.lock.Unlock()
	m.lock.Lock()
	result := make([]*Job, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}
	return result
}

type JobService interface {
	GetJob(id string) (*Job, error)
	SetJob(id string, job *Job) error
}

type NopJobService struct{}

func (s NopJobService) GetJob(id string) (*Job, error) {
	return nil, nil
}

func (s NopJobService) SetJob(id string, job *Job) error {
	return nil
}
