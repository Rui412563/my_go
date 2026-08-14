package logic

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"greet/internal/svc"
	"greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
	// 记录请求参数
	l.Infof("收到 Greet 请求, Name: %s", req.Name)

	resp = &types.Response{}
	switch req.Name {
	case "you":
		resp.Message = "Hello to you ,hello word，你好!"
	case "me":
		resp.Message = "Hello to me,hello word，你好!"
	default:
		// 未匹配到任何分支时记录警告日志
		l.Infof("未匹配到 Name: %s, 返回空响应", req.Name)
	}

	l.Infof("Greet 响应: %s", resp.Message)
	return
}

// 新增：视频播放逻辑
// 注意：这个方法的签名和普通 API 不同，因为它需要直接操作 http.ResponseWriter
func (l *GreetLogic) PlayVideo(req *types.VideoRequest, w http.ResponseWriter, r *http.Request) error {
	// 1. 拼接视频文件路径
	// 假设 video 文件夹和 main.go 在同一级目录
	videoPath := filepath.Join("video", filepath.Base(req.VideoName))
	l.Infof("请求播放视频, VideoName: %s, 完整路径: %s", req.VideoName, videoPath)

	// 2. 打开文件
	f, err := os.Open(videoPath)
	if err != nil {
		// 文件不存在，返回 404
		l.Errorf("视频文件打开失败, 路径: %s, 错误: %v", videoPath, err)
		http.Error(w, "Video not found", http.StatusNotFound)
		return nil
	}
	defer f.Close()

	// 3. 获取文件信息
	fi, err := f.Stat()
	if err != nil {
		l.Errorf("获取视频文件信息失败, 路径: %s, 错误: %v", videoPath, err)
		http.Error(w, "Cannot stat video", http.StatusInternalServerError)
		return nil
	}

	// 4. 使用 ServeContent 发送文件
	// 它会自动处理 Content-Type, Content-Length, 以及支持拖动进度条的 Range 请求
	l.Infof("开始发送视频, 文件名: %s, 大小: %d bytes", fi.Name(), fi.Size())
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	l.Infof("视频发送完成, 文件名: %s", fi.Name())

	return nil
}
