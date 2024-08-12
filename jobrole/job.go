package jobrole

import (
	"sort"
	"sync"
)

type Job struct {
	ID       string
	Title    string
	Desc     string
	DutyList []string
}

func (j *Job) MergeID(id string) *Job {
	j.ID = id
	return j
}
func (j *Job) MergeTitle(title string) *Job {
	j.Title = title
	return j
}
func (j *Job) MergeDesc(desc string) *Job {
	j.Desc = desc
	return j
}
func (j *Job) AppendDuty(dutyid string) *Job {
	j.DutyList = append(j.DutyList, dutyid)
	return j
}
func NewJob() *Job {
	return &Job{
		DutyList: []string{},
	}
}

type Jobs []*Job

func (jobs Jobs) Len() int {
	return len(jobs)
}
func (jobs Jobs) Less(i, j int) bool {
	return jobs[i].ID < jobs[j].ID
}
func (jobs Jobs) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}

type JobMap struct {
	NopJobService
	data map[string]*Job
	lock sync.Mutex
}

func NewJobMap() *JobMap {
	return &JobMap{
		data: map[string]*Job{},
	}
}

func (m *JobMap) GetJob(id string) (*Job, error) {
	defer m.lock.Unlock()
	m.lock.Lock()
	return m.data[id], nil
}

func (m *JobMap) SetJob(job *Job) {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.data[job.ID] = job
}
func (m *JobMap) List() []*Job {
	defer m.lock.Unlock()
	m.lock.Lock()
	result := make([]*Job, 0, len(m.data))
	for _, v := range m.data {
		result = append(result, v)
	}
	sort.Sort(Jobs(result))
	return result
}

type JobService interface {
	GetJob(id string) (*Job, error)
}

type NopJobService struct{}

func (s NopJobService) GetJob(id string) (*Job, error) {
	return nil, nil
}
