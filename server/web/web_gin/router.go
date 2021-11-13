package web_gin

import (
	"net/http"
	"time"

	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/middleware"
	gomonitor "github.com/szuecs/gin-gomonitor"
	ginmon "github.com/szuecs/gin-gomonitor/aspects"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

type CustomAspect struct {
	CustomValue int
}

func (a *CustomAspect) GetStats() interface{} {
	return a.CustomValue
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}

// InitRouter 初始化模块路由
func (ws *WebServer) InitRouter() error {
	ws.app.Use(limit.MaxAllowed(50))
	router := ws.app.Group("/")
	{
		router.Use(middleware.Cors()) // 如需跨域可以打开

		requestAspect := ginmon.NewRequestTimeAspect()
		requestAspect.StartTimer(5 * time.Second)

		counterAspect := ginmon.NewCounterAspect()
		counterAspect.StartTimer(3 * time.Second)

		genericAspect := ginmon.NewGenericChannelAspect("generic")
		genericAspect.StartTimer(3 * time.Second)
		genericCH := genericAspect.SetupGenericChannelAspect()

		asps := []aspects.Aspect{counterAspect, requestAspect, genericAspect}
		router.Use(ginmon.CounterHandler(counterAspect))
		// curl http://localhost:9000/RequestTime
		router.Use(ginmon.RequestTimeHandler(requestAspect))
		// curl http://localhost:9000/
		gomonitor.Start(9000, asps)
		// last middleware
		router.Use(gin.Recovery())

		// each request to all handlers like below will increment the Counter
		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"Counter": map[string]string{
					"msg": "Request Counter - Loook at http://localhost:9000/Counter",
					"cmd": "curl http://localhost:9000/Counter ; for i in {1..20}; do curl localhost:8080/ &>/dev/null ; curl localhost:8080/foo &>/dev/null ; done; sleep 3; curl http://localhost:9000/Counter"},
				"RequestTime": map[string]string{
					"msg": "RequestTime is registered at http://localhost:9000/RequestTime and will return data after 5 seconds.",
					"cmd": "for j in {0..100}; do for i in {1..20}; do curl localhost:8080/ ; done; sleep 0.5; curl localhost:9000/RequestTime ; done"},
				"GenericChannelAspect": map[string]string{
					"msg": "Generic Aspect can process arbitrary map[string]float64 data - Loook at http://localhost:9000/generic",
					"cmd": "curl http://localhost:9000/generic ; for i in {1..20}; do curl localhost:8080/generic &>/dev/null ; done; sleep 3; curl http://localhost:9000/generic"}})
		})

		router.GET("/generic", func(ctx *gin.Context) {
			for i := 0; i < 100; i++ {
				genericCH <- ginmon.DataChannel{Name: "foo", Value: float64(i % 2)}
				genericCH <- ginmon.DataChannel{Name: "bar", Value: float64(i % 5)}
			}
			ctx.JSON(http.StatusOK, gin.H{
				"GenericChannelAspect": map[string]string{
					"msg": "Generic Aspect can process arbitrary map[string]float64 data - Loook at http://localhost:9000/generic",
					"cmd": "curl http://localhost:9000/generic ; for i in {1..20}; do curl localhost:8080/generic &>/dev/null ; done; sleep 3; curl http://localhost:9000/generic"}})
		})
		// for _, party := range ws.parties {
		// 	app.PartyFunc(party.Perfix, party.PartyFunc)
		// }
	}
	if ws.staticPrefix != "" {
		ws.AddUploadStatic()
	}
	// if ws.webPrefix != "" {
	// 	ws.AddWebStatic()
	// }
	return nil
}

// GetSources 获取系统路由
// - PermRoutes 权鉴路由
// - NoPermRoutes 公共路由
func (ws *WebServer) GetSources() ([]map[string]string, []map[string]string) {
	routeLen := len(ws.app.Routes())
	permRoutes := make([]map[string]string, 0, routeLen)
	noPermRoutes := make([]map[string]string, 0, routeLen)
	for _, r := range ws.app.Routes() {

		route := map[string]string{
			"path": r.Path,
			"name": r.Handler,
			"act":  r.Method,
		}
		// handerNames := context.HandlersNames(r.HandlerFunc)
		// if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) || !arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
		// 	noPermRoutes = append(noPermRoutes, route)
		// } else {
		permRoutes = append(permRoutes, route)
		// }
	}
	return permRoutes, noPermRoutes
}
