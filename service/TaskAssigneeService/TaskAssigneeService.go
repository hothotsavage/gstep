package TaskAssigneeService

import (
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/util/db/dao"
	"gorm.io/gorm"
)

func ToVO(taskAssignee entity.TaskAssignee, tx *gorm.DB) vo.TaskAssigneeVO {
	aVO := vo.TaskAssigneeVO{}
	aVO.TaskAssignee = taskAssignee
	user := dao.GetById[entity.User](taskAssignee.UserId, tx)
	aVO.User = *user

	return aVO
}
