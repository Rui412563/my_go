# 第一阶段：编译
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
# 设置国内 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .

# 编译为 Linux 可执行文件，关闭 CGO 生成纯静态二进制
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o greet greet.go

# 第二阶段：运行（极简基础镜像）
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为北京时间
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /root/
# 从第一阶段拷贝编译好的二进制文件
COPY --from=builder /app/greet .
# 拷贝 go-zero 的配置文件
COPY --from=builder /app/etc ./etc
# 拷贝视频文件夹
COPY video ./video

# 创建日志目录
RUN mkdir -p /root/logs

# 声明数据卷
VOLUME ["/root/logs"]

EXPOSE 8888
CMD ["./greet", "-f", "etc/greet-api.yaml"]