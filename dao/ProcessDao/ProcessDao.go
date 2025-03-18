package ProcessDao

import (
	"fmt"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"gorm.io/gorm"
)

func Id2Entity(id int, tx *gorm.DB) *entity.User {
	var user entity.User
	err := tx.Raw("select * from process where userid=? limit 1", id).First(&user).Error
	if nil != err {
		msg := fmt.Sprintf("查询流程实例(id=%s)失败: %s", id, err)
		panic(ServerError.New(msg))
	}
	return &user
}
