package MouldHandler

import (
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/dao/MouldDao"
	"github.com/hothotsavage/gstep/dao/TemplateDao"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/MouldService"
	"github.com/hothotsavage/gstep/util/db/dao"
	PageDto "github.com/hothotsavage/gstep/util/db/page"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Save(writer http.ResponseWriter, request *http.Request) {
	mould := entity.Mould{}
	RequestParsUtil.Body2dto(request, &mould)

	tx := ctx.GetTx(request)

	MouldService.Save(&mould, tx)

	AjaxJson.SuccessByData(mould.Id).Response(writer)
}

func List(writer http.ResponseWriter, request *http.Request) {
	pageDto := PageDto.PageDto{}
	RequestParsUtil.Body2dto(request, &pageDto)

	tx := ctx.GetTx(request)
	pagination := MouldDao.Query(pageDto, tx)
	AjaxJson.SuccessByData(pagination).Response(writer)
}

func Detail(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DetailDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	entity := dao.CheckById[entity.Mould](dto.Id, tx)
	AjaxJson.SuccessByData(entity).Response(writer)
}

func Delete(writer http.ResponseWriter, request *http.Request) {
	dto := dto.DeleteDto{}
	RequestParsUtil.Body2dto(request, &dto)

	tx := ctx.GetTx(request)

	cnt := TemplateDao.TemplateCount(dto.Id, tx)
	if cnt > 0 {
		AjaxJson.Fail("请先删除流程图").Response(writer)
		return
	}

	dao.DeleteById[entity.Mould](dto.Id, tx)

	AjaxJson.Success().Response(writer)
}
