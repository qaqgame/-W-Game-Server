package Parser

import (
	"wGame/Model"
	"encoding/json"
	"fmt"
)


func CreateRes(data Model.Res) string {
	resp, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(resp)
}

func CreateResErr(data Model.ResErr) string {
	resp, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(resp)
}