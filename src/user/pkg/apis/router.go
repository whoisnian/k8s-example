package apis

type Message string

const (
	MsgSuccess              Message = "success"
	MsgInternalError        Message = "internal error"
	MsgDisableRegistration  Message = "disable registration"
	MsgInvalidCookie        Message = "invalid cookie"
	MsgInvalidParams        Message = "invalid params"
	MsgInvalidEmail         Message = "invalid email"
	MsgEmailOrPasswordError Message = "email or password error"
)

type Code int

const (
	CodeSuccess         Code = 0
	CodeInternalError   Code = 1001
	CodeInvalidCookie   Code = 1002
	CodeInvalidParams   Code = 1003
	CodeInvalidEmail    Code = 1004
	CodeInvalidPassword Code = 1005
)

type MessageResponse struct {
	Message Message `json:"message"`
}

// https://opentelemetry.io/docs/specs/semconv/http/http-spans/#status
//
// For HTTP status codes in the 4xx range span status MUST be left unset in case of SpanKind.SERVER and MUST be set to Error in case of SpanKind.CLIENT.
// So we use custom Code to identify internal response status.
type InternalResponse[T any] struct {
	Code    Code    `json:"code"`
	Message Message `json:"message"`
	Data    T       `json:"data"`
}

func NewInternalResponse[T any](code Code, message Message, data T) *InternalResponse[T] {
	return &InternalResponse[T]{code, message, data}
}
