package verbose

import (
	"log"
	"os"
)

var Loger *log.Logger

func InitLoger(path string) {
	file := path + "\\CSO2-Server.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	Loger = log.New(logFile, "", log.LstdFlags|log.LUTC) // 将文件设置为loger作为输出
}
