package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/NidzamuddinMuzakki/movies-abishar/common"
	constants "github.com/NidzamuddinMuzakki/movies-abishar/common/constant"
)

type Response struct {
	HTTPStatus   int         `json:"-"`
	Code         int         `json:"code"`
	Success      bool        `json:"success"`
	RequestID    string      `json:"request_id,omitempty"`
	MessageTitle string      `json:"messageTitle,omitempty"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	Error        any         `json:"error,omitempty"`
}

type ResponseSnakeCase struct {
	Success      bool        `json:"success"`
	MessageTitle string      `json:"message_title,omitempty"`
	Message      string      `json:"message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type WithPagination struct {
	CurrentPage  int64       `json:"current_page"`
	NextPage     int64       `json:"next_page"`
	PreviousPage int64       `json:"previous_page"`
	PerPage      int64       `json:"per_page"`
	Total        int64       `json:"total"`
	TotalPage    int64       `json:"total_page"`
	Data         interface{} `json:"data"`
}

type DetailError struct {
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	ErrorDetail      any    `json:"errorDetail,omitempty"`
}

const (
	CodeSuccess             = 200
	CodeInternalServerError = 500
	CodeUnimplemented       = 501
	CodeBadGateway          = 502
	CodeServiceUnavailable  = 503
	CodeGatewayTimeout      = 504
	CodeBadRequest          = 400
	CodeUnauthorized        = 401
	CodeForbidden           = 403
	CodeNotFound            = 404
	CodeConflict            = 409
	CodeUnprocessableEntity = 422
	CodeTooManyRequest      = 429

	MessageTitleSuccess      = "Successfully!"
	MessageTitleBadRequest   = "Bad Request!"
	MessageTitleUnauthorized = "Unauthorized!"
	MessageTitleForbidden    = "Action Forbidden!"

	MessageSuccess             = "success to fetch data"
	MessageDependencyFailed    = "dependency failed"
	MessageUnimplemented       = "method unimplemented"
	MessageInternalServerError = "oops! something wrong"
	MessageInvalidArgument     = "data invalid"
	MessageUnauthorized        = "invalid access token"
	MessageForbidden           = "access token don't have permission to access this action"
	MessageConflict            = "action conflict"
	MessageNotFound            = "data not found"
	MessageUnprocessableEntity = "unprocessable entity"
	MessageTooManyRequest      = "too many request"
	MessageGatewayTimeout      = "gateway timeout"
	MessageServiceUnavailable  = "service unavailable"
	MessageBadRequest          = "Invalid data request"
	MessageInvalidSignature    = "Invalid Signature"

	ErrorCodeDataNotFound = "DATA_NOT_FOUND"
	ErrorCodeInvalid      = "INVALID_REQUEST"
)

// OK define response with status OK
func OK(ctx context.Context, data interface{}) *Response {
	return &Response{
		HTTPStatus:   http.StatusOK,
		Success:      true,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeSuccess,
		MessageTitle: MessageTitleSuccess,
		Message:      MessageSuccess,
		Data:         data,
	}
}

func CUSTOMMESSAGEOK(ctx context.Context, data interface{}, messageTittle string, mesage string) *Response {
	return &Response{
		HTTPStatus:   http.StatusOK,
		Success:      true,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeSuccess,
		MessageTitle: messageTittle,
		Message:      mesage,
		Data:         data,
	}
}

// Created define response with status created
func Created(ctx context.Context, data interface{}) *Response {
	if data == nil {
		data = "success"
	}
	return &Response{
		HTTPStatus: http.StatusCreated,
		Success:    true,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeSuccess,
		Data:       data,
	}
}

// InvalidArgument define response invalid argument from client
func InvalidArgument(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusBadRequest,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeBadRequest,
		Message:    MessageInvalidArgument,
	}
}

// BadGateway define response for dependency service failed, such as DB, cache, other upstream service problem
func BadGateway(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusBadGateway,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeBadGateway,
		Message:    MessageDependencyFailed,
	}
}

func BadRequest(ctx context.Context) *Response {
	return &Response{
		HTTPStatus:   http.StatusBadRequest,
		Success:      false,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeBadRequest,
		MessageTitle: MessageTitleBadRequest,
		Message:      MessageBadRequest,
	}
}

func BadRequestWithMessage(ctx context.Context, message string) *Response {
	return &Response{
		HTTPStatus:   http.StatusBadRequest,
		Success:      false,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeBadRequest,
		MessageTitle: MessageTitleBadRequest,
		Message:      message,
	}
}

// Unimplemented ...
func Unimplemented(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusNotImplemented,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeUnimplemented,
		Message:    MessageUnimplemented,
	}
}

// InternalServerError define response for unknown internal server error
// Please avoid return this if we can categorize it to more specific such as BadGateway
func InternalServerError(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusInternalServerError,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeInternalServerError,
		Message:    MessageInternalServerError,
	}
}

// Forbidden define response if client request data that is not belong for their access
func Forbidden(ctx context.Context) *Response {
	return &Response{
		HTTPStatus:   http.StatusForbidden,
		Success:      false,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeForbidden,
		Message:      MessageForbidden,
		MessageTitle: MessageTitleForbidden,
	}
}

// UnprocessableEntity define response unprocessable entity from client
func UnprocessableEntity(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusUnprocessableEntity,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeUnprocessableEntity,
		Message:    MessageUnprocessableEntity,
	}
}

// Conflict define response conflict action from client
func Conflict(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusConflict,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeConflict,
		Message:    MessageConflict,
	}
}

// NotFound define response for requested resource not found
func NotFound(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusNotFound,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeNotFound,
		Message:    MessageNotFound,
	}
}

// TooManyRequest define response that request is too many
func TooManyRequest(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusTooManyRequests,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeTooManyRequest,
		Message:    MessageTooManyRequest,
	}
}

// Unauthorised ...
func Unauthorised(ctx context.Context) *Response {
	return &Response{
		HTTPStatus:   http.StatusUnauthorized,
		Success:      false,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeUnauthorized,
		Message:      MessageUnauthorized,
		MessageTitle: MessageTitleUnauthorized,
	}
}

// Signature not valid ...
func InvalidSignature(ctx context.Context) *Response {
	return &Response{
		HTTPStatus:   http.StatusBadRequest,
		Success:      false,
		RequestID:    common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:         CodeBadRequest,
		Message:      MessageInvalidSignature,
		MessageTitle: MessageTitleBadRequest,
	}
}

// GatewayTimeout define response timeout
func GatewayTimeout(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusGatewayTimeout,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeGatewayTimeout,
		Message:    MessageGatewayTimeout,
	}
}

// Unavailable define response for service unavailable
func Unavailable(ctx context.Context) *Response {
	return &Response{
		HTTPStatus: http.StatusServiceUnavailable,
		Success:    false,
		RequestID:  common.GetContextValueAsString(ctx, constants.XRequestIDHeader),
		Code:       CodeServiceUnavailable,
		Message:    MessageServiceUnavailable,
	}
}

// WithRequestID modifies api response's requestID
// no need to call this unless you have a specific requirement to change the requestID
func (r *Response) WithRequestID(requestID string) *Response {
	r.RequestID = requestID
	return r
}

// WithCode modifies api response's code
func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

// WithMessage modifies api response's message
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// WithMessageTitle modifies api response's message
func (r *Response) WithMessageTitle(message string) *Response {
	r.MessageTitle = message
	return r
}

// ToMap modifies api response's message
func (r *Response) ToMap() map[string]interface{} {
	var respMap map[string]interface{}
	byteResp, _ := json.Marshal(r)
	err := json.Unmarshal(byteResp, &respMap)
	if err != nil {
		// TODO log error here
		return nil
	}

	return respMap
}

// ToHTTPCodeAndMap to convert response to return http status + map[string] interface
func (r *Response) ToHTTPCodeAndMap() (int, map[string]interface{}) {
	return r.HTTPStatus, r.ToMap()
}

func CustomErrorResponse(errorCode string, errorMessage string, errorDetail interface{}) Response {
	res := Response{
		Code:    CodeBadRequest,
		Success: false,
		Message: MessageBadRequest,
		Error: DetailError{
			ErrorCode:        errorCode,
			ErrorDescription: errorMessage,
			ErrorDetail:      errorDetail,
		},
	}
	return res
}
