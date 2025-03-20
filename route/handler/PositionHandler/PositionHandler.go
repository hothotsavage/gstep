package PositionHandler

import (
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/dao/PositionDao"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"net/http"
)

func GetPositions(writer http.ResponseWriter, request *http.Request) {
	tx := ctx.GetTx(request)
	positions := PositionDao.GetPositions(tx)
	AjaxJson.SuccessByData(positions).Response(writer)
}
