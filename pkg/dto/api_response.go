package dto

import (
	"errors"
	"strconv"
)

type StatusMessage int64

const (
	Success StatusMessage = iota
	Fail
	Error
)

func (r StatusMessage) MarshalJSON() ([]byte, error) {

	switch r {
	case Success:
		return []byte(strconv.Quote("success")), nil
	case Fail:
		return []byte(strconv.Quote("fail")), nil
	case Error:
		return []byte(strconv.Quote("error")), nil
	}

	return nil, errors.New("unsupported response status")

}

type ResponseObject struct {
	StatusMsg StatusMessage `json:"status"` //(success|fail|error)
	Message   string        `json:"message,omitempty"`
	Data      interface{}   `json:"data,omitempty"`
}

func (r *ResponseObject) SetStatus(status StatusMessage) {
	r.StatusMsg = status
}
