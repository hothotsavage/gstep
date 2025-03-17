package ProcessHandler

import (
	"github.com/hothotsavage/gstep/model/dto"
	"github.com/hothotsavage/gstep/model/entity"
	"github.com/hothotsavage/gstep/service/ProcessService"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/db/dao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"net/http"
)

func Start(writer http.ResponseWriter, request *http.Request) {
	processStartDto := dto.ProcessStartDto{}
	RequestParsUtil.Body2dto(request, &processStartDto)

	tx := DbUtil.GetTx()
	dao.CheckById[entity.User](processStartDto.UserId, tx)
	//创建流程及启动任务
	vo := ProcessService.Start(&processStartDto, tx)

	//任务状态变更通知
	//TaskService.NotifyTasksStateChange(id, tx)

	tx.Commit()

	AjaxJson.SuccessByData(vo).Response(writer)
}

func Pass(writer http.ResponseWriter, request *http.Request) {
	processPassDto := dto.ProcessPassDto{}
	RequestParsUtil.Body2dto(request, &processPassDto)

	tx := DbUtil.GetTx()
	dao.CheckById[entity.User](processPassDto.UserId, tx)
	//审核通过
	vo := ProcessService.Pass(processPassDto, tx)

	tx.Commit()

	AjaxJson.SuccessByData(vo).Response(writer)
}

// 退回到指定上一步
func Refuse(writer http.ResponseWriter, request *http.Request) {
	processRefuseDto := dto.ProcessRefuseDto{}
	RequestParsUtil.Body2dto(request, &processRefuseDto)

	tx := DbUtil.GetTx()
	dao.CheckById[entity.User](processRefuseDto.UserId, tx)
	//拒绝
	vo := ProcessService.Refuse(processRefuseDto, tx)

	tx.Commit()

	AjaxJson.SuccessByData(vo).Response(writer)
}
