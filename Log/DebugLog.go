package Log

import (
	"os"
	"fmt"
	"log"
)

func DebugLog(fileName, debuginfo string) {
	logFile, err := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
	if err != nil {
		fmt.Println("open file err:",err)
	}
	defer logFile.Close()

	debuglogger := log.New(logFile,"[Debug]",log.Ldate|log.Lmicroseconds|log.Lshortfile)
	debuglogger.Println(debuginfo)
}
