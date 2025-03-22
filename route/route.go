package route

import (
	"fmt"
	"github.com/hothotsavage/gstep/config"
	"github.com/hothotsavage/gstep/ctx"
	"github.com/hothotsavage/gstep/route/handler/DepartmentHandler"
	"github.com/hothotsavage/gstep/route/handler/MouldHandler"
	"github.com/hothotsavage/gstep/route/handler/NotifyHandler"
	"github.com/hothotsavage/gstep/route/handler/PositionHandler"
	"github.com/hothotsavage/gstep/route/handler/ProcessHandler"
	"github.com/hothotsavage/gstep/route/handler/TaskHandler"
	"github.com/hothotsavage/gstep/route/handler/TemplateHandler"
	"github.com/hothotsavage/gstep/route/handler/TestHandler"
	"github.com/hothotsavage/gstep/util/ServerError"
	"github.com/hothotsavage/gstep/util/db/DbUtil"
	"github.com/hothotsavage/gstep/util/net/AjaxJson"
	"github.com/hothotsavage/gstep/util/net/RequestParsUtil"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

var Mux = http.NewServeMux()

func middleware(h http.HandlerFunc) http.HandlerFunc {
	handler := authHandle(h)
	handler = transaction(handler)
	handler = jsonResponseHead(handler)
	handler = crossOrigin(handler)
	return handler
}

func noAuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	handler := transaction(h)
	handler = jsonResponseHead(handler)
	//本地调试时需要处理跨域
	if config.Config.IsDebugLocal {
		//接入spring cloud gateway后,不需要处理跨域
		//spring cloud gateway 已经处理跨域
		handler = crossOrigin(handler)
	}
	return handler
}

func jsonResponseHead(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		h(w, r)
	}
}

func authHandle(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := config.Config.Auth.Secret
		token := RequestParsUtil.GetAuthorizationToken(r)

		if secret != token {
			panic(ServerError.New("无访问权限"))
		}

		h(w, r)
	}
}

//func errorHandle(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		defer func() {
//			err := recover()
//
//			tx := DbUtil.GetTx()
//			tx.Rollback()
//
//			if nil != err {
//				debug.PrintStack()
//				AjaxJson.Fail(fmt.Sprintf("%s", err)).Response(w)
//			}
//		}()
//
//		h(w, r)
//	}
//}

func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-requested-with,Authorization")
		//w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		//注意:Access-Control-Allow-Origin不能设置成*
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		if len(r.Header.Get("Origin")) > 0 {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		}
		//if len(r.Header.Get("Referer")) > 0 {
		//	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Referer"))
		//}

		//预请求只返回响应头
		if "OPTIONS" == r.Method {
			//注意:w.WriteHeader(http.StatusAccepted)之后的w.Header().Set代码无效
			w.WriteHeader(http.StatusAccepted)
			return
		}

		h(w, r)
	}
}

// 事务及错误处理
func transaction(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := DbUtil.Db.Begin()

		defer func() {
			if er := recover(); er != nil {
				tx.Rollback()

				log.Println(er)
				log.Println(string(debug.Stack()))

				AjaxJson.Fail(fmt.Sprintf("%s", er)).Response(w)
				return
			}
		}()

		//将事务对象写入context
		r = ctx.SetTx(r, tx)

		h(w, r)
		tx.Commit()
	}
}

func Setup() {
	setupRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Config.Port),
		Handler:      Mux,
		ReadTimeout:  time.Duration(30 * int(time.Second)),
		WriteTimeout: time.Duration(30 * int(time.Second)),
		//MaxHeaderBytes: 1 << 20,
	}
	log.Printf("web server start up at port%s", server.Addr)
	err := server.ListenAndServe()

	if nil != err {
		log.Printf("server startup fail: %v", err)
	}
}

// define route
func setupRoutes() {
	//1.流程模板
	//保存
	Mux.HandleFunc("/mould/save", noAuthMiddleware(MouldHandler.Save))
	//查询
	Mux.HandleFunc("/mould/list", noAuthMiddleware(MouldHandler.List))
	//详情
	Mux.HandleFunc("/mould/detail", noAuthMiddleware(MouldHandler.Detail))
	Mux.HandleFunc("/mould/delete", noAuthMiddleware(MouldHandler.Delete))

	Mux.HandleFunc("/template/new_draft", noAuthMiddleware(TemplateHandler.NewDraft))
	Mux.HandleFunc("/template/release", noAuthMiddleware(TemplateHandler.Release))
	//保存
	Mux.HandleFunc("/template/save", noAuthMiddleware(TemplateHandler.Save))
	//查询
	Mux.HandleFunc("/template/query", noAuthMiddleware(TemplateHandler.Query))
	//详情
	Mux.HandleFunc("/template/detail", noAuthMiddleware(TemplateHandler.Detail))
	//删除
	Mux.HandleFunc("/template/delete", noAuthMiddleware(TemplateHandler.Delete))

	//2.流程实例
	//启动流程
	Mux.HandleFunc("/process/start", noAuthMiddleware(ProcessHandler.Start))
	//任务审核
	Mux.HandleFunc("/process/pass", noAuthMiddleware(ProcessHandler.Pass))
	//任务驳回
	Mux.HandleFunc("/process/refuse", noAuthMiddleware(ProcessHandler.Refuse))
	//查询驳回步骤列表
	Mux.HandleFunc("/process/refuse_prevsteps", noAuthMiddleware(ProcessHandler.RefusePrevSteps))
	//详情
	Mux.HandleFunc("/process/detail", noAuthMiddleware(ProcessHandler.Detail))

	//查询我的任务
	Mux.HandleFunc("/task/pending", noAuthMiddleware(TaskHandler.Pending))
	Mux.HandleFunc("/task/query", noAuthMiddleware(TaskHandler.Query))

	//+++ 测试接口 ++++++++++++++++++
	//接收通知
	Mux.HandleFunc("/notify/task_state_change", noAuthMiddleware(NotifyHandler.TaskStateChange))

	//部门查询
	Mux.HandleFunc("/department/get_child_department", noAuthMiddleware(DepartmentHandler.GetChildDepartments))
	Mux.HandleFunc("/department/get_users", noAuthMiddleware(DepartmentHandler.GetUsers))

	//职位查询
	Mux.HandleFunc("/position/positions", noAuthMiddleware(PositionHandler.GetPositions))

	//职位查询
	Mux.HandleFunc("/test/hello", noAuthMiddleware(TestHandler.Hello))
}
