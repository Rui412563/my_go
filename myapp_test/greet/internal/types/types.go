package types

type Request struct {
	Name string `path:"name,options=you|me"`
}

type Response struct {
	Message string `json:"message"`
}

// 新增：视频播放请求结构体
type VideoRequest struct {
	VideoName string `path:"videoName"`
}
