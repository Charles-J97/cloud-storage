package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// RespMsg : http响应数据的通用结构
type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// JSONBytes : 对象转json格式的二进制数组
func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return r
}

// JSONString : 对象转json格式的string
func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return string(r)
}

//统一转发和接收GET和DELETE请求
func HandleGETAndDELRequest(url, method string) (*RespMsg, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respMsg := &RespMsg{}
	err = json.Unmarshal(body, respMsg)
	if err != nil {
		return nil, err
	}
	return respMsg, nil
}

//统一转发和接收POST和PUT请求
func HandlePOSTAndPUTRequest(jsonBody []byte, url, method string) (*RespMsg, error) {
	reader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	//必须设定该参数,POST参数才能正常提交，意思是以json串提交数据
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respMsg := &RespMsg{}
	err = json.Unmarshal(body, respMsg)
	if err != nil {
		return nil, err
	}

	return respMsg, nil
}
