package errors

import (
	"encoding/json"
	"time"
)

type RequeueError struct {
	errMsg          string
	requeueDuration time.Duration
}

func (r *RequeueError) Error() string {
	return r.errMsg
}

func (r *RequeueError) Duration() time.Duration {
	return r.requeueDuration
}

func NewRequeueError(msg string, requeueSeconds int) error {
	return &RequeueError{
		errMsg:          msg,
		requeueDuration: time.Second * time.Duration(requeueSeconds),
	}
}

type APIError struct {
	ErrMsg string `json:"error"`
}

func (e *APIError) Error() string { return e.ErrMsg }

func (e *APIError) ErrorJSON() []byte {
	out, _ := json.MarshalIndent(e, "  ", "")
	return append(out, []byte("\n")...)
}

func NewAPIError(msg string) error {
	return &APIError{ErrMsg: msg}
}

func IsRequeueError(err error) (*RequeueError, bool) {
	if qerr, ok := err.(*RequeueError); ok {
		return qerr, true
	}
	return nil, false
}

func IsAPIError(err error) (*APIError, bool) {
	if qerr, ok := err.(*APIError); ok {
		return qerr, true
	}
	return nil, false
}
