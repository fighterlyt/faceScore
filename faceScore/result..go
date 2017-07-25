package faceScore

import "errors"

const (
	successCode = "success"
	failCode    = "fail"
)

type Result struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Result  float32 `json:"result"`
}

func (r Result) IsSuccess() bool {
	return r.Code == successCode
}

func (r Result) GetSocre() float32 {
	return r.Result
}

func (r Result) HasError() error {
	if r.Code == failCode {
		return errors.New(r.Message)
	}
	return nil
}
