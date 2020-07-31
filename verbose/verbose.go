package verbose

import (
	"log"
)

var (
	Level     uint32
	LogFile   uint32
	IsConsole uint32
)

func DebugInfo(l uint32, v ...interface{}) {
	if l <= Level {
		if IsConsole == 0 {
			log.Println(v...)
		}
		if LogFile != 0 {
			Loger.Println(v...)
		}
	}
}
