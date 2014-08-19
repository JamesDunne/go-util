package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseKind int

const (
	Undetermined ResponseKind = iota
	HTML
	JSON
	Empty
)

type Error struct {
	ResponseKind ResponseKind
	StatusCode   int
	Error        error
}

// Override this to provide custom web error logging:
type ErrorLogFunc func(req *http.Request, werr *Error)

func defaultErrorLog(req *http.Request, werr *Error) {
	if werr == nil {
		return
	}
	// Don't log non-error HTTP statuses:
	if werr.StatusCode < 400 {
		return
	}

	err := ""
	if werr.Error != nil {
		err = werr.Error.Error()
	}

	log.Printf("%3d %s %s ERROR %s\n", werr.StatusCode, req.Method, req.URL, err)
}

var DefaultErrorLog ErrorLogFunc = ErrorLogFunc(defaultErrorLog)

func NewError(err error, statusCode int, kind ResponseKind) *Error {
	return &Error{
		ResponseKind: kind,
		StatusCode:   statusCode,
		Error:        err,
	}
}

func AsError(err error, statusCode int) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		ResponseKind: Undetermined,
		StatusCode:   statusCode,
		Error:        err,
	}
}

func AsErrorHTML(err error, statusCode int) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		ResponseKind: HTML,
		StatusCode:   statusCode,
		Error:        err,
	}
}

func AsErrorJSON(err error, statusCode int) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		ResponseKind: JSON,
		StatusCode:   statusCode,
		Error:        err,
	}
}

func (e *Error) Respond(rsp http.ResponseWriter) bool {
	if e == nil {
		return false
	}

	if e.ResponseKind == HTML {
		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")
		rsp.WriteHeader(e.StatusCode)
		rsp.Write([]byte(e.Error.Error()))
		return true
	} else if e.ResponseKind == JSON {
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
	} else if e.ResponseKind == Empty {
		rsp.WriteHeader(e.StatusCode)
		return true
	}

	return false
}

func (e *Error) AsJSON() *Error {
	if e == nil {
		return nil
	}
	e.ResponseKind = JSON
	return e
}

func (e *Error) AsHTML() *Error {
	if e == nil {
		return nil
	}
	e.ResponseKind = HTML
	return e
}

func (e *Error) As(kind ResponseKind) *Error {
	if e == nil {
		return nil
	}
	e.ResponseKind = kind
	return e
}

type ErrorHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) *Error
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) *Error

// ServeHTTP calls f(w, r).
func (f ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) *Error {
	return f(w, r)
}

func Log(logfunc ErrorLogFunc, h ErrorHandler) ErrorHandler {
	if logfunc == nil {
		logfunc = DefaultErrorLog
	}
	return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) (werr *Error) {
		werr = h.ServeHTTP(w, r)
		if werr != nil {
			logfunc(r, werr)
		}
		return werr
	})
}

func ReportErrors(h ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		werr := h.ServeHTTP(w, r)
		werr.Respond(w)
		return
	})
}

///////////////////////////////////////////////////////////

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

	return AsErrorJSON(err, statusCode).Respond(rsp)
}
