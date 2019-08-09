package Log

import (
	"wGame/Global"
	"runtime"
	"path/filepath"
	"strings"
	"strconv"
)

func LogController() {
	for true {
		select {
		case loginfo := <-Global.DebugLogger:
			DebugLog(Global.LogFileName,loginfo)
		}
	}
}

func GetTransferInfo() string {
	funcname := ""
	pc, filename, line, ok := runtime.Caller(2)
	if ok {
		funcname = runtime.FuncForPC(pc).Name()       // main.(*MyStruct).foo
		funcname = filepath.Ext(funcname)             // .foo
		funcname = strings.TrimPrefix(funcname, ".")  // foo

		filename = filepath.Base(filename)  // /full/path/basename.go => basename.go
	}
	posinfo := filename+" funcname:"+funcname+" line:"+strconv.Itoa(line)+" "
	return posinfo
}