package TemplateDao

import (
	"fmt"
	"github.com/hothotsavage/gstep/enum/TemplateState"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/ServerError"
	"gorm.io/gorm"
)

func GetLatestVersion(mouldId int, tx *gorm.DB) *entity.Template {
	var entities []*entity.Template
	err := tx.Where("mould_id=? and state=?", mouldId, TemplateState.RELEASE.Code).Order("version desc").Find(&entities).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("未找到流程模板(mould_id=%d) %s", mouldId, err)))
	}

	if len(entities) == 0 {
		return nil
	}

	return entities[0]
}

func GetTemplate(mouldId int, versionId int, tx *gorm.DB) *entity.Template {
	var entities []*entity.Template
	err := tx.Where("mould_id=? and version=?", mouldId, versionId).Order("version desc").Find(&entities).Error
	if nil != err {
		panic(err)
	}

	if len(entities) == 0 {
		return nil
	}

	return entities[0]
}

func NewTemplateId(tx *gorm.DB) int {
	maxMouldId := 0
	err := tx.Raw("select ifnull(max(mould_id),0) from template").Scan(&maxMouldId).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("获取新mouldId失败,%v", err)))
	}

	return maxMouldId + 1
}

func NewVersion(mouldId int, tx *gorm.DB) int {
	maxVersion := 0
	err := tx.Raw("select ifnull(max(version),0) from template where mould_id=?", mouldId).Scan(&maxVersion).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询模板(mouldId=%d)最大版本号失败,%s", mouldId, err)))
	}

	return maxVersion + 1
}

func TemplateCount(mouldId int, tx *gorm.DB) int {
	cnt := 0
	err := tx.Raw("select count(*) from template where mould_id=?", mouldId).Scan(&cnt).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程图(mouldId=%d)数量失败,%s", mouldId, err)))
	}

	return cnt
}
