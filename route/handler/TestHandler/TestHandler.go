package TestHandler

import (
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"net/http"
)

func Hello(writer http.ResponseWriter, request *http.Request) {
	b := 0
	r := 1 / b
	print(r)
	AjaxJson.Success().Response(writer)
}
