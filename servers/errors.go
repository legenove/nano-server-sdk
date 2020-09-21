package servers

import (
	"github.com/legenove/utils"
)

type ServerError struct {
	statusCode int      `json:"-"`
	Code       string   `json:"code"`
	Msg        string   `json:"msg"`
	Details    []string `json:"details"`
}

func (sErr *ServerError) Error() string {
	if len(sErr.Details) == 0 {
		return sErr.Msg
	}
	out := make([]string, len(sErr.Details)+2)
	out[0] = sErr.Msg
	out[1] = " : "
	copy(out[2:], sErr.Details[:])
	return utils.ConcatenateStrings(out...)
}

var ServerErrorMap = make(map[string]*ServerError)

func (sErr *ServerError) New(details []string, error_code ...string) *ServerError {
	var code string
	if len(error_code) > 0 {
		code = error_code[0]
	} else {
		code = sErr.Code
	}
	return &ServerError{
		statusCode: sErr.statusCode,
		Code:       code,
		Msg:        sErr.Msg,
		Details:    details,
	}
}

func (sErr *ServerError) StatusCode() int {
	if sErr.statusCode > 0 {
		return sErr.statusCode
	}
	return 200
}

func (sErr *ServerError) SetStatusCode(code int) *ServerError {
	sErr.statusCode = code
	return sErr
}

func NewServerError(msg string, code string, statusCode int) *ServerError {
	apiErr := &ServerError{
		statusCode: statusCode,
		Code:       code,
		Msg:        msg,
		Details:    []string{},
	}
	ServerErrorMap[msg] = apiErr
	return apiErr
}

var (
	ErrUnKnowRequest         = NewServerError("unknow_error", "10000", 400)
	ErrProjectValidator      = NewServerError("project_validator_error", "10002", 400)
	ErrProjectMatch          = NewServerError("project_match_error", "10003", 400)
	ErrPageNotFoundRequest   = NewServerError("not_found", "10004", 404)
	ErrMethodNotAllowRequest = NewServerError("no_method", "10005", 405)
	ErrSchemaOptionNotFound  = NewServerError("doc_not_found", "10006", 404)
	ErrUnDefineRequest       = NewServerError("undefined_error", "10007", 400)
	ErrRequestErr            = NewServerError("requests_error", "10008", 400)
	ErrGetRequestHost        = NewServerError("get_request_host_error", "10009", 400)
)
