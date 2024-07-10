package authErrs

type errorCode string

const (
	ENOTFOUND  errorCode = "user_not_found"
	EEXIST     errorCode = "user_already_exists"
	EINCORRECT errorCode = "incorrect_data"
	EINTERNAL  errorCode = "internal"
)

type Error struct {
	Code    errorCode `json:"code"`
	Message string    `json:"message"`
	Op      string    `json:"op"`
}

func New(code errorCode, msg, op string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
		Op:      op,
	}
}
