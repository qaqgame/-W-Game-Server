package Parser

import (
	"wGame/Model"
	"encoding/json"
	"fmt"
)


func CreateRes(data Model.Res) []byte {
	resp, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
		return nil
	}
	return resp
}

func CreateResErr(data Model.ResErr) []byte {
	resp, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
		return nil
	}
	return resp
}

func CreateReconndata(data Model.StateSync) []byte {
	resp, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}