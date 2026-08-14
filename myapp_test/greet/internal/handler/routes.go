package handler

import (
	"net/http"

	"greet/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/from/:name",
				Handler: GreetHandler(serverCtx),
			},
			// 新增：注册视频播放路由
			{
				Method:  http.MethodGet,
				Path:    "/video/:videoName",
				Handler: PlayVideoHandler(serverCtx),
			},
		},
	)
}
