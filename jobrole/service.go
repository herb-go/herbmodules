package jobrole

import (
	"time"

	"github.com/herb-go/herbsecurity/authorize/role"
	"github.com/herb-go/herbsecurity/authorize/role/roleparser"
)

type Service struct {
	DutyService
	JobService
	AssignmentService
}

func (s *Service) LoadAssignmentsRoles(assignments []*Assignment) ([]*role.Roles, error) {
	result := []*role.Roles{}
	for _, v := range assignments {
		job, err := s.JobService.GetJob(v.ID)
		if err != nil {
			return nil, err
		}
		for _, duty := range job.DutyList {
			roles, err := roleparser.Parse(duty.Roles)
			if err != nil {
				return nil, err
			}
			result = append(result, roles)
		}
	}
	return result, nil
}
func (s *Service) GetUserRoles(uid string, t *time.Time) ([]*role.Roles, error) {
	assignments, err := s.AssignmentService.GetUserAssignments(uid)
	if err != nil {
		return nil, err
	}

	var timestamp int64
	if t == nil {
		timestamp = time.Now().Unix()
	} else {
		timestamp = t.Unix()
	}
	result := make([]*Assignment, 0, len(assignments))
	for _, v := range assignments {
		if v.Before != nil && *(v.Before) < timestamp {
			continue
		}
		if v.After != nil && *(v.After) > timestamp {
			continue
		}
		result = append(result, v)
	}
	return s.LoadAssignmentsRoles(result)
}
func New() *Service {
	return &Service{
		DutyService:       NopDutyServive{},
		JobService:        NopJobService{},
		AssignmentService: NopAssignmentService{},
	}
}
