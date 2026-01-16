package mGin

// CustomError custom error
type CustomError struct {
	HTTPCode    int
	Code        int
	Message     string
	ErrorInfo   interface{}
	DeclineCode string
}

func (cErr CustomError) Error() string {
	return cErr.Message
}

func (cErr CustomError) WithDeclineCode(code string) CustomError {
	return CustomError{
		HTTPCode:    cErr.HTTPCode,
		Code:        cErr.Code,
		Message:     cErr.Message,
		ErrorInfo:   cErr.ErrorInfo,
		DeclineCode: code,
	}
}

type Errors struct {
	Meta Meta      `json:"meta"`
	Data *struct{} `json:"data"`
}
