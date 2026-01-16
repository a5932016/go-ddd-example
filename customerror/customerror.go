package customerror

import (
	"fmt"
	"net/http"

	"github.com/a5932016/go-ddd-example/util/mGin"
)

var (
	// auth
	AccountNotFound = mGin.CustomError{
		HTTPCode: http.StatusNotFound,
		Code:     10001,
		Message:  "Account not found",
	}
	WrongPassword = mGin.CustomError{
		HTTPCode: http.StatusUnauthorized,
		Code:     10002,
		Message:  "Wrong password",
	}
	InvalidSession = mGin.CustomError{
		HTTPCode: http.StatusUnauthorized,
		Code:     10003,
		Message:  "Invalid session",
	}
	NoPermission = mGin.CustomError{
		HTTPCode: http.StatusUnauthorized,
		Code:     10004,
		Message:  "No permission",
	}
	PasswordTooLong = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10005,
		Message:  "Password too long",
	}
	InvalidPassword = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10006,
		Message:  "Invalid password",
	}
	InvalidPasswordConfirmation = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10007,
		Message:  "Invalid password confirmation",
	}
	InvalidAccount = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10008,
		Message:  "Invalid account",
	}
	InvalidUserName = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10009,
		Message:  "Invalid user name",
	}
	DuplicateUserAccount = mGin.CustomError{
		HTTPCode: http.StatusConflict,
		Code:     10010,
		Message:  "User account has been used",
	}
	NoHierarchyPermission = mGin.CustomError{
		HTTPCode: http.StatusUnauthorized,
		Code:     10011,
		Message:  "No hierarchy permission",
	}
	NoSelfUpdatePermission = mGin.CustomError{
		HTTPCode: http.StatusUnauthorized,
		Code:     10012,
		Message:  "No self update permission",
	}
	InvalidEmail = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     10013,
		Message:  "Invalid email",
	}
	// functions
	RecordNotFound = mGin.CustomError{
		HTTPCode: http.StatusNotFound,
		Code:     30001,
		Message:  "Record not found",
	}
	InvalidForeignId = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30002,
		Message:  "Invalid foreign Id",
	}
	InvalidName = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     30003,
		Message:  "Only allow 'a-z', 'A-Z', '0-9', '_', '-', '(', and ')' characters",
	}
	DuplicateName = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     30004,
		Message:  "Duplicate name",
	}
	DuplicateTitle = mGin.CustomError{
		HTTPCode: http.StatusNotAcceptable,
		Code:     30005,
		Message:  "Duplicate title",
	}
	InvalidUserId = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30006,
		Message:  "Invalid user Id",
	}
	InvalidImageExtension = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30007,
		Message:  "Invalid image extension",
	}
	DuplicateImage = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30008,
		Message:  "Duplicate image",
	}
	InvalidImageSize = func(bits int64) mGin.CustomError {
		return mGin.CustomError{
			HTTPCode: http.StatusNotAcceptable,
			Code:     30009,
			Message:  fmt.Sprintf("Each image file size should be under %.0fKB", BitsToKB(bits)),
		}
	}
	InvalidImageContent = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30010,
		Message:  "Invalid image content",
	}
	InvalidYouTubeURL = mGin.CustomError{
		HTTPCode: http.StatusBadRequest,
		Code:     30011,
		Message:  "Invalid YouTube URL",
	}
	SomeRecordsNotFound = mGin.CustomError{
		HTTPCode: http.StatusNotFound,
		Code:     30012,
		Message:  "Some records not found",
	}
)

func BitsToKB(bits int64) float64 {
	const bitsPerKB = 8 * 1024
	return float64(bits) / bitsPerKB
}
