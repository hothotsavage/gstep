package ProcessHandler

import (
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/ProcessService"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Start(writer http.ResponseWriter, request *http.Request) {
	processStartDto := dto.ProcessStartDto{}
	RequestParsUtil.Body2dto(request, &processStartDto)
	tx := ctx.GetTx(request)

	dao.CheckById[entity.User](processStartDto.UserId, tx)
	//创建流程及启动任务
	vo := ProcessService.Start(&processStartDto, tx)

	AjaxJson.SuccessByData(vo).Response(writer)
}

func Pass(writer http.ResponseWriter, request *http.Request) {
	processPassDto := dto.ProcessPassDto{}
	RequestParsUtil.Body2dto(request, &processPassDto)

	tx := ctx.GetTx(request)

	dao.CheckById[entity.User](processPassDto.UserId, tx)
	//审核通过
	vo := ProcessService.Pass(processPassDto, tx)

	AjaxJson.SuccessByData(vo).Response(writer)
}

// 退回到指定上一步
func Refuse(writer http.ResponseWriter, request *http.Request) {
	processRefuseDto := dto.ProcessRefuseDto{}
	RequestParsUtil.Body2dto(request, &processRefuseDto)

	tx := ctx.GetTx(request)
	dao.CheckById[entity.User](processRefuseDto.UserId, tx)
	//拒绝
	vo := ProcessService.Refuse(processRefuseDto, tx)

	AjaxJson.SuccessByData(vo).Response(writer)
}

// 获取驳回到的之前任务列表
func RefusePrevSteps(writer http.ResponseWriter, request *http.Request) {
	refusePrevStepsDTO := dto.RefusePrevStepsDTO{}
	RequestParsUtil.Body2dto(request, &refusePrevStepsDTO)

	tx := ctx.GetTx(request)
	dao.CheckById[entity.Process](refusePrevStepsDTO.ProcessId, tx)
	//拒绝
	steps := ProcessService.RefusePrevSteps(refusePrevStepsDTO.ProcessId, tx)

	AjaxJson.SuccessByData(steps).Response(writer)
}

func Detail(writer http.ResponseWriter, request *http.Request) {
	detailDTO := dto.DetailDto{}
	RequestParsUtil.Body2dto(request, &detailDTO)

	tx := ctx.GetTx(request)
	process := dao.CheckById[entity.Process](detailDTO.Id, tx)
	vo := ProcessService.ToDetailVO(process, tx)

	AjaxJson.SuccessByData(vo).Response(writer)
}
