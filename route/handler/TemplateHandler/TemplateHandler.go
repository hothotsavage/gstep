package TemplateHandler

import (
	"fmt"
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/dao/ProcessDao"
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/enum/TemplateState"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/TemplateService"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func NewDraft(writer http.ResponseWriter, request *http.Request) {
	tx := ctx.GetTx(request)

	dto := dto.TemplateNewDto{}
	RequestParsUtil.Body2dto(request, &dto)

	mould := dao.CheckById[entity.Mould](dto.MouldId, tx)

	//克隆最新版本作为草稿
	latestTemplate := TemplateDao.GetLatestVersion(dto.MouldId, tx)
	if latestTemplate != nil {
		latestTemplate.Id = 0
		latestTemplate.State = TemplateState.DRAFT.Code
		latestTemplate.Version = latestTemplate.Version + 1
		AjaxJson.SuccessByData(latestTemplate).Response(writer)
		return
	}

	//未找到最新版本,克隆默认模板
	defaultMouldId := 1
	defaultTemplate := TemplateDao.GetLatestVersion(defaultMouldId, tx)
	if defaultTemplate == nil {
		AjaxJson.Fail("未找到默认模板").Response(writer)
		return
	}

	defaultTemplate.MouldId = dto.MouldId
	defaultTemplate.Id = 0
	defaultTemplate.Title = mould.Title
	defaultTemplate.State = TemplateState.DRAFT.Code
	defaultTemplate.Version = 1
	AjaxJson.SuccessByData(defaultTemplate).Response(writer)
}

func Save(writer http.ResponseWriter, request *http.Request) {
	newTemplate := entity.Template{}
	RequestParsUtil.Body2dto(request, &newTemplate)

	tx := ctx.GetTx(request)

	dao.CheckById[entity.Mould](newTemplate.MouldId, tx)

	//已存在版本
	if newTemplate.Id > 0 {
		dao.CheckById[entity.Template](newTemplate.Id, tx)
		dao.SaveOrUpdate(&newTemplate, tx)
	} else { //已存在模板,新版本
		version := TemplateDao.NewVersion(newTemplate.MouldId, tx)
		newTemplate.Version = version
		newTemplate.State = TemplateState.DRAFT.Code
		dao.SaveOrUpdate(&newTemplate, tx)
	}

	AjaxJson.SuccessByData(newTemplate.Id).Response(writer)
}

func Release(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DetailDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	pEntity := dao.CheckById[entity.Template](dto.Id, tx)
	if pEntity.State == TemplateState.RELEASE.Code {
		AjaxJson.Fail("重复发布").Response(writer)
		return
	}
	pEntity.State = TemplateState.RELEASE.Code
	dao.SaveOrUpdate(pEntity, tx)
	AjaxJson.SuccessByData(*pEntity).Response(writer)
}

func Query(writer http.ResponseWriter, request *http.Request) {
	dto := dto.TemplateQueryDto{}
	RequestParsUtil.Body2dto(request, &dto)
	tx := ctx.GetTx(request)

	pagintaion := TemplateService.Query(&dto, tx)
	AjaxJson.SuccessByData(pagintaion).Response(writer)
}

func Detail(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DetailDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	entity := dao.CheckById[entity.Template](dto.Id, tx)
	AjaxJson.SuccessByData(entity).Response(writer)
}

func Delete(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DeleteDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	cnt := ProcessDao.TemplateProcessCount(dto.Id, tx)
	if cnt > 0 {
		AjaxJson.Fail(fmt.Sprintf("流程图已实例化%d次,无法删除", cnt)).Response(writer)
		return
	}

	dao.DeleteById[entity.Template](dto.Id, tx)

	AjaxJson.Success().Response(writer)
}
