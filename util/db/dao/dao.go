package dao

import (
	"fmt"
	"github.com/hothotsavage/gstep/util/LocalTime"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/entity"
	"gorm.io/gorm"
	"reflect"
	"time"
)

// 插入或更新记录
func SaveOrUpdate(pEntity any, tx *gorm.DB) {
	//反射获取id
	value := reflect.ValueOf(pEntity)
	id := reflect.Indirect(value).FieldByName("Id").Interface()

	//反射写入更新时间
	if reflect.Indirect(value).FieldByName("UpdatedAt").IsValid() {
		time := LocalTime.LocalTime(time.Now())
		now := reflect.ValueOf(&time)
		reflect.Indirect(value).FieldByName("UpdatedAt").Set(now)
	}

	//根据id类型判断插入还是更新记录
	switch id.(type) {
	case int:
		if id.(int) < 1 {
			result := tx.Create(pEntity)
			if nil != result.Error {
				panic(result.Error)
			}
		} else {
			tx.Save(pEntity)
		}
	case string:
		if len(id.(string)) < 1 {
			result := tx.Create(pEntity)
			if nil != result.Error {
				panic(result.Error)
			}
		} else {
			tx.Save(pEntity)
		}
	}
}

func GetById[T entity.CommonEntity, I int | string](id I, tx *gorm.DB) *T {
	var detail T

	err := tx.Table(detail.TableName()).Where("id=?", id).First(&detail).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("未找到(表:%s id=%s)记录:%s", detail.TableName(), id, err.Error())))
	}
	return &detail
}

func CheckById[T entity.CommonEntity, I int | string](id I, tx *gorm.DB) *T {
	var detail T
	err := tx.Table(detail.TableName()).Where("id=?", id).First(&detail).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("未找到(表:%s id=%d)记录:%s", detail.TableName(), id, err.Error())))
	}

	newId := detail.GetId()
	err = CheckId(newId)
	if nil != err {
		panic(err)
	}

	return &detail
}

func CheckId(id any) error {
	switch id.(type) {
	case int:
		if id == 0 {
			return ServerError.New("无效的表id")
		}
	case string:
		if len(id.(string)) == 0 {
			return ServerError.New("无效的表id")
		}
	}

	return nil
}

func DeleteById[T entity.CommonEntity, I int | string](id I, tx *gorm.DB) {
	var detail T
	res := tx.Delete(&detail, id)
	if res.Error != nil {
		fmt.Println("删除(table=%s)失败 ", detail.TableName(), res.Error)
	}
}
