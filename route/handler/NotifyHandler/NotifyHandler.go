package NotifyHandler

import (
	"github.com/hothotsavage/gstep/model/vo"
	"github.com/hothotsavage/gstep/util/JsonUtil"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"log"
	"net/http"
)

func TaskStateChange(writer http.ResponseWriter, request *http.Request) {
	dto := vo.TaskStateChangeNotifyVo{}
	RequestParsUtil.Body2dto(request, &dto)

	log.Println("接收到任务变更通知消息")
	log.Println(JsonUtil.Obj2PrettyJson(dto))

	AjaxJson.Success().Response(writer)
}
