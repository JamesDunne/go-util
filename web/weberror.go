package web

import (
	"encoding/json"
	"net/http"
)

type WebError struct {
	StatusCode int
	Error      error
}

func AsWebError(err error, statusCode int) *WebError {
	if err == nil {
		return nil
	}
	return &WebError{
		StatusCode: statusCode,
		Error:      err,
	}
}

func (e *WebError) Do(action func(e *WebError)) *WebError {
	if e != nil {
		action(e)
	}
	return e
}

func (e *WebError) RespondHTML(rsp http.ResponseWriter) bool {
	if e == nil {
		return false
	}

	rsp.WriteHeader(e.StatusCode)
	rsp.Write([]byte(e.Error.Error()))
	return true
}

func (e *WebError) RespondJSON(rsp http.ResponseWriter) bool {
	if e == nil {
		return false
	}

	rsp.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp.WriteHeader(e.StatusCode)
	j, jerr := json.Marshal(&struct {
		StatusCode int    `json:"statusCode"`
		Error      string `json:"error"`
	}{
		StatusCode: e.StatusCode,
		Error:      e.Error.Error(),
	})
	if jerr != nil {
		panic(jerr)
	}
	rsp.Write(j)
	return true
}

func JsonSuccess(rsp http.ResponseWriter, result interface{}) {
	rsp.Header().Set("Content-Type", "application/json; charset=utf-8")
	rsp.WriteHeader(http.StatusOK)
	j, jerr := json.Marshal(&struct {
		StatusCode int         `json:"statusCode"`
		Result     interface{} `json:"result"`
	}{
		StatusCode: http.StatusOK,
		Result:     result,
	})
	if jerr != nil {
		panic(jerr)
	}
	rsp.Write(j)
}

func JsonErrorIf(rsp http.ResponseWriter, err error, statusCode int) bool {
	if err == nil {
		return false
	}

	return AsWebError(err, statusCode).RespondJSON(rsp)
}
