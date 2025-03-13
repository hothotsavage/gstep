package LogUtil

import (
	"github.com/hothotsavage/gstep/util/JsonUtil"
	"log"
)

func PrintPretty(obj any) {
	jsonStr := JsonUtil.Obj2PrettyJson(obj)
	log.Printf("%s", jsonStr)
}
