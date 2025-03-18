FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache libc6-compat

# 设置工作目录
WORKDIR /app

# 复制main文件
COPY main /app/main
COPY ./static /app/static
COPY ./templates /app/templates

# 赋予执行权限
RUN chmod +x /app/main

# 设置入口点
ENTRYPOINT ["./main"]

# 暴露端口
EXPOSE 80
