package UserDao

import (
	"fmt"
	"github.com/hothotsavage/gstep/dao/DepartmentDao"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"gorm.io/gorm"
)

func GetDepartmentUsers(departmentId string, tx *gorm.DB) []entity.User {
	users := []entity.User{}

	err := tx.Raw("select * from user a "+
		" where a.department_id=?"+
		" order by a.id asc ", departmentId).Scan(&users).Error
	if nil != err {
		msg := fmt.Sprintf("找不到部门员工: %s", err)
		panic(ServerError.New(msg))
	}
	return users
}

func GetDepartmentUserCount(departmentId string, tx *gorm.DB) int {
	cnt := 0

	err := tx.Raw("select count(1) from user a "+
		" where a.department_id=?"+
		" order by a.id asc ", departmentId).Scan(&cnt).Error
	if nil != err {
		msg := fmt.Sprintf("找不到部门员工数量: %s", err)
		panic(ServerError.New(msg))
	}
	return cnt
}

func GetGrandsonDepartmentUsers(departmentId string, tx *gorm.DB) []entity.User {
	users := []entity.User{}

	departmentIds := DepartmentDao.GetGrandsonDepartmentIds(departmentId, tx)

	err := tx.Raw("select * from user a "+
		" where a.department_id in ?"+
		" order by a.id asc ", departmentIds).Scan(&users).Error
	if nil != err {
		msg := fmt.Sprintf("找不到部门员工: %s", err)
		panic(ServerError.New(msg))
	}
	return users
}

func IsUserInDepartment(userId string, departmentId string, tx *gorm.DB) bool {
	cnt := 0

	err := tx.Raw("select count(1) from user a "+
		" where a.department_id=?"+
		" and a.id=? ", departmentId, userId).Scan(&cnt).Error
	if nil != err {
		msg := fmt.Sprintf("找不到部门员工: %s", err)
		panic(ServerError.New(msg))
	}
	return cnt > 0
}

func IsUserInDepartments(userId string, departments []entity.Department, tx *gorm.DB) bool {
	if len(departments) == 0 {
		return false
	}

	for _, v := range departments {
		isIn := IsUserInDepartment(userId, v.Id, tx)
		if isIn {
			return true
		}
	}

	return false
}
