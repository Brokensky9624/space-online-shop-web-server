package response

import "space.online.shop.web.server/util/tool"

const (
	success = "success"
	fail    = "fail"
)

func NewResp() *RespObj {
	return &RespObj{
		Status: success,
		Data:   make([]interface{}, 0),
	}
}

type RespObj struct {
	Status  string        `json:"status"`
	ErrMsg  string        `json:"errMsg"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

func (resp *RespObj) SetStatus(status string) *RespObj {
	resp.Status = status
	return resp
}

func (resp *RespObj) SetErrMsg(errMsg string) *RespObj {
	resp.ErrMsg = errMsg
	return resp
}

func (resp *RespObj) SetMessage(message string) *RespObj {
	resp.Message = message
	return resp
}

func (resp *RespObj) SetData(dataList ...interface{}) *RespObj {
	resp.Data = append(resp.Data, dataList...)
	return resp
}

func SuccessRespObj(message string, dataList ...interface{}) *RespObj {
	resp := NewResp().SetMessage(message)
	if dataList != nil {
		dataList = tool.FilterSlice(dataList, func(data interface{}) bool {
			return data != nil
		})
		if len(dataList) > 0 {
			resp.SetData(dataList...)
		}
	}
	return resp
}

func FailRespObj(err error) *RespObj {
	return NewResp().SetStatus(fail).SetErrMsg(err.Error())
}
