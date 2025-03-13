package vo

import "github.com/hothotsavage/gstep/model/entity"

type CandidateVo struct {
	entity.Candidate
	Department entity.Department
	User       entity.User
	Position   string
}
