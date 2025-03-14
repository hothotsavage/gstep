package PositionHandler

import (
	"github.com/hothotsavage/gstep/dao/PositionDao"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"net/http"
)

func GetPositions(writer http.ResponseWriter, request *http.Request) {
	positions := PositionDao.GetPositions(DbUtil.Db)
	AjaxJson.SuccessByData(positions).Response(writer)
}
