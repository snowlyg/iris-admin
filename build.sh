#!/bin/bash
# 停止容器
sudo docker stop irisadminapi_demo
# 删除镜像
sudo docker rmi -f irisadminapi:demo

PWD=$(pwd)
cd $PWD/www
# 打包前端文件
sudo npm run build:stage
cd ../
#编译前端文件，如果在配置文件开启 bindata 模式
go generate
# 此处为 mac 跨平台编译参数，如果是其他平台需要相应修改
CGO_ENABLED=1 GOOS=linux CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc go build -a -installsuffix cgo -ldflags "-w -s -X main.Version=v0.0.1" -o ./main_lin
# 构建镜像
sudo docker build -t irisadminapi:demo .
# 启动容器
docker run -d -p 8085:8085  --name irisadminapi_demo irisadminapi:demo

exit 0