package TemplateService

import (
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/util/db/dao"
	"gorm.io/gorm"
	"strconv"
)

func Query(dto *dto.TemplateQueryDto, tx *gorm.DB) []entity.Template {
	list := []entity.Template{}
	sql := "select id,template_id,title,version,created_at,updated_at,deleted_at from template " +
		" where 1=1 "
	if dto.VersionId > 0 {
		sql = sql + " and id=" + strconv.Itoa(dto.VersionId)
	}
	sql = sql + " limit " + strconv.Itoa(dto.Limit)
	sql = sql + " offset " + strconv.Itoa((dto.Page-1)*dto.Limit)

	err := tx.Raw(sql).Scan(&list).Error
	if nil != err {
		panic(err)
	}

	return list
}

func QueryDetail(dto *dto.TemplateQueryDetailDto, tx *gorm.DB) *entity.Template {
	pTemplate := &entity.Template{}
	if dto.Id > 0 {
		pTemplate = dao.CheckById[entity.Template](dto.Id, tx)
		return pTemplate
	} else if dto.TemplateId > 0 && dto.VersionId > 0 {
		pTemplate = TemplateDao.GetTemplate(dto.TemplateId, dto.VersionId, tx)
	} else if dto.TemplateId > 0 {
		pTemplate = TemplateDao.GetLatestVersionByTemplateId(dto.TemplateId, tx)
	}

	return pTemplate
}

func QueryInfo(dto *dto.TemplateQueryInfoDto, tx *gorm.DB) *entity.Template {
	pTemplate := &entity.Template{}
	if dto.VersionId > 0 {
		pTemplate = dao.CheckById[entity.Template](dto.VersionId, tx)
	} else if dto.TemplateId > 0 {
		pTemplate = TemplateDao.GetLatestVersionByTemplateId(dto.TemplateId, tx)
	}

	pTemplate.RootStep = entity.Step{}

	return pTemplate
}
