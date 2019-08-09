package Parser

import (
	"wGame/Model"
	"encoding/json"
	"wGame/Log"
	"wGame/Global"
)


func CreateRes(data Model.Res) string {
	resp, err := json.Marshal(&data)
	if err != nil {
		//fmt.Println(err)
		loginfo := Log.GetTransferInfo()
		Global.DebugLogger <- loginfo + err.Error()
		return ""
	}
	return string(resp)
}

func CreateResErr(data Model.ResErr) string {
	resp, err := json.Marshal(&data)
	if err != nil {
		//fmt.Println(err)
		loginfo := Log.GetTransferInfo()
		Global.DebugLogger <- loginfo + err.Error()
		return ""
	}
	return string(resp)
}