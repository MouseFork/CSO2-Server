package verbose

import (
	"log"
	"os"
)

var Loger *log.Logger

func InitLoger(path, name string) {
	file := path + "\\" + name
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		panic(err)
	}
	Loger = log.New(logFile, "", log.LstdFlags|log.LUTC) // 将文件设置为loger作为输出
}
