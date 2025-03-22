package TemplateService

import (
	"fmt"
	"github.com/hothotsavage/gstep/dao/ProcessDao"
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/dao"
	pagination "github.com/hothotsavage/gstep/util/net/page"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"strconv"
)

func Query(dto *dto.TemplateQueryDto, tx *gorm.DB) *pagination.Pagination[vo.TemplateVO] {
	list := []entity.Template{}
	baseSql := "select * from template " +
		" where 1=1 "
	if dto.MouldId > 0 {
		baseSql = baseSql + " and mould_id=" + strconv.Itoa(dto.MouldId)
	}
	if dto.VersionId > 0 {
		baseSql = baseSql + " and version=" + strconv.Itoa(dto.VersionId)
	}
	baseSql = baseSql + " order by version desc "
	pageSql := baseSql + " limit " + strconv.Itoa(dto.Limit)
	pageSql = pageSql + " offset " + strconv.Itoa((dto.Page-1)*dto.Limit)

	err := tx.Raw(pageSql).Scan(&list).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程图列表失败 %s", err)))
	}

	total := 0
	countSql := "select count(*) " +
		" from (" + baseSql + ") t "
	err = tx.Raw(countSql).Scan(&total).Error
	if nil != err {
		panic(ServerError.New(fmt.Sprintf("查询流程图总数失败 %s", err)))
	}

	res := pagination.Pagination[vo.TemplateVO]{}
	vos := lo.Map(list, func(detail entity.Template, index int) vo.TemplateVO {
		return *ToVO(&detail, tx)
	})
	res.List = vos
	res.Total = total
	return &res
}

func QueryInfo(dto *dto.TemplateQueryInfoDto, tx *gorm.DB) *entity.Template {
	pTemplate := &entity.Template{}
	if dto.VersionId > 0 {
		pTemplate = dao.CheckById[entity.Template](dto.VersionId, tx)
	} else if dto.TemplateId > 0 {
		pTemplate = TemplateDao.GetLatestVersion(dto.TemplateId, tx)
	}

	pTemplate.RootStep = entity.Step{}

	return pTemplate
}

func ToVO(template *entity.Template, tx *gorm.DB) *vo.TemplateVO {
	vo := &vo.TemplateVO{}
	vo.Template = *template

	cnt := ProcessDao.TemplateProcessCount(template.Id, tx)
	vo.ProcessCount = cnt

	return vo
}
