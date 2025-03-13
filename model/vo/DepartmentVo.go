package vo

import "github.com/hothotsavage/gstep/model/entity"

type DepartmentVo struct {
	entity.Department
	HasSubDepartments bool `json:"hasSubDepartments"`
	UserCount         int  `json:"userCount"`
}
