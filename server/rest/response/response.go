package response

type RespObj struct {
	Status string        `json:"status"`
	ErrMsg string        `json:"errMsg"`
	Data   []interface{} `json:"data"`
}

const (
	success = "success"
	fail    = "fail"
)

func SuccessRespObj(data []interface{}) *RespObj {
	respObj := &RespObj{
		Status: success,
		ErrMsg: "",
		Data:   make([]interface{}, 0),
	}
	respObj.Data = append(respObj.Data, data...)
	return respObj
}

func FailRespObj(err error) *RespObj {
	respObj := &RespObj{
		Status: fail,
		ErrMsg: err.Error(),
		Data:   make([]interface{}, 0),
	}
	return respObj
}
