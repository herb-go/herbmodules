package jobrole

import (
	"strings"

	"github.com/herb-go/herbsecurity/authorize/role"
	"github.com/herb-go/herbsecurity/authorize/role/roleparser"
)

type RolesLoader interface {
	LoadRoles(s *Service, a *Assignment) ([]*role.Roles, error)
}

type PlainRolesLoader struct{}

func (l PlainRolesLoader) LoadRoles(s *Service, a *Assignment) ([]*role.Roles, error) {
	result := []*role.Roles{}
	if a == nil {
		return nil, nil
	}
	job, err := s.JobService.GetJob(a.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}
	for _, dutyid := range job.DutyList {
		d, err := s.DutyService.GetDuty(dutyid)
		if err != nil {
			return nil, err
		}
		if d != nil {
			roles, err := roleparser.Parse(d.Roles)
			if err != nil {
				return nil, err
			}
			result = append(result, roles)
		}
	}
	return result, nil
}

type ReplacerRolesLoader struct{}

func (l ReplacerRolesLoader) LoadRoles(s *Service, a *Assignment) ([]*role.Roles, error) {
	result := []*role.Roles{}
	if a == nil {
		return nil, nil
	}
	job, err := s.JobService.GetJob(a.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}
	var replacer *strings.Replacer
	if len(a.Data) > 0 {
		replacertoken := []string{}
		for k, v := range a.Data {
			replacertoken = append(replacertoken, k, v)
		}
		replacer = strings.NewReplacer(replacertoken...)
	}
	for _, dutyid := range job.DutyList {
		d, err := s.DutyService.GetDuty(dutyid)
		if err != nil {
			return nil, err
		}
		if d != nil {
			rolesstr := d.Roles
			if replacer != nil {
				rolesstr = replacer.Replace(rolesstr)
			}
			roles, err := roleparser.Parse(rolesstr)
			if err != nil {
				return nil, err
			}
			result = append(result, roles)
		}
	}
	return result, nil

}
