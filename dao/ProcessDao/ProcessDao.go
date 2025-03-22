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

func TemplateProcessCount(templateId int, tx *gorm.DB) int {
	cnt := 0
	err := tx.Raw("select count(*) from process where template_id=?", templateId).Scan(&cnt).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程实例(templateId=%d)数量失败,%s", templateId, err)))
	}

	return cnt
}

func MouldProcessCount(mouldId int, tx *gorm.DB) int {
	cnt := 0
	err := tx.Raw("select count(*) from process "+
		" where exists(select 1 from template t "+
		" where t.id=process.template_id "+
		" and t.mould_id=?)", mouldId).Scan(&cnt).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程实例(mouldId=%d)数量失败,%s", mouldId, err)))
	}

	return cnt
}
