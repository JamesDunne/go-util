package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type WebError struct {
	StatusCode int
	Error      error
}

// Override this to provide custom web error logging:
type WebErrorLogFunc func(req *http.Request, werr *WebError)

var DefaultWebErrorLog WebErrorLogFunc = WebErrorLogFunc(defaultWebErrorLog)

func defaultWebErrorLog(req *http.Request, werr *WebError) {
	log.Printf("%3d %s %s ERROR %s\n", werr.StatusCode, req.Method, req.URL, werr.Error.Error())
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

type WebErrorHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) *WebError
}

type WebErrorHandlerFunc func(http.ResponseWriter, *http.Request) *WebError

// ServeHTTP calls f(w, r).
func (f WebErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) *WebError {
	return f(w, r)
}

func Log(logfunc WebErrorLogFunc, h WebErrorHandler) WebErrorHandler {
	if logfunc == nil {
		logfunc = DefaultWebErrorLog
	}
	return WebErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) (werr *WebError) {
		werr = h.ServeHTTP(w, r)
		if werr != nil {
			logfunc(r, werr)
		}
		return werr
	})
}

func JSON(h WebErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		werr := h.ServeHTTP(w, r)
		if werr != nil {
			werr.RespondJSON(w)
		}
		return
	})
}

func HTML(h WebErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		werr := h.ServeHTTP(w, r)
		if werr != nil {
			werr.RespondHTML(w)
		}
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

	return AsWebError(err, statusCode).RespondJSON(rsp)
}
