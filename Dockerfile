#镜像
FROM golang:1.14.2 as build
ENV GO111MODULE off
#设置工作目录
WORKDIR $GOPATH/src/github.com/KouKouChan/CSO2-Server
#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/github.com/KouKouChan/CSO2-Server
COPY ./kerlong $GOPATH/src/github.com/KouKouChan/CSO2-Server/kerlong
#go构建可执行文件
RUN go build .
#暴露端口
EXPOSE 30001
EXPOSE 30002
#最终运行docker的命令
ENTRYPOINT  ["./CSO2-Server"]