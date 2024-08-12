package jobrole

import (
	"time"

	"github.com/herb-go/herbsecurity/authorize/role"
)

type Service struct {
	DutyService
	JobService
	AssignmentService
	RolesLoader
}

func (s *Service) LoadAssignmentsRoles(assignments []*Assignment) ([]*role.Roles, error) {
	result := []*role.Roles{}
	for _, assignment := range assignments {
		roles, err := s.RolesLoader.LoadRoles(s, assignment)
		if err != nil {
			return nil, err
		}
		result = append(result, roles...)
	}
	return result, nil
}
func (s *Service) GetUserRoles(uid string, t *int64) (*role.Roles, error) {
	assignments, err := s.AssignmentService.GetUserAssignments(uid)
	if err != nil {
		return nil, err
	}

	var timestamp int64
	if t == nil {
		timestamp = time.Now().Unix()
	} else {
		timestamp = *t
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
	userroles, err := s.LoadAssignmentsRoles(result)

	if err != nil {
		return nil, err
	}
	return role.Concat(userroles...), nil
}
func New() *Service {
	return &Service{
		DutyService:       NopDutyService{},
		JobService:        NopJobService{},
		AssignmentService: NopAssignmentService{},
		RolesLoader:       PlainRolesLoader{},
	}
}
