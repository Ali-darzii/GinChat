package validators

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"strings"
	"unicode"
)

func PhoneNoValidate(fl validator.FieldLevel) bool {
	phoneNo := fl.Field().String()
	for _, item := range phoneNo {
		if !unicode.IsDigit(item) {
			return false
		}
	}
	if !strings.HasPrefix(phoneNo, "09") {
		return false
	}
	return true
}
func UsernameValidate(fl validator.FieldLevel) bool {
	// handling ->  not first not last == underScore
	username := fl.Field().String()
	const underScore uint8 = 95

	if username[0] == underScore || username[len(username)-1] == underScore {
		return false
	}
	for index, item := range username {
		// handling -> _ == 95 -> not together
		if item == int32(underScore) && username[index+1] == underScore {
			return false
		}
		// only digits and letter and _
		if !unicode.IsDigit(item) && !unicode.IsLetter(item) && item != int32(underScore) {
			return false
		}

	}

	return false
}
func NameValidate(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name != "" {
		return len(name) > 3
	}
	return true
}
func ImageValidate(fl validator.FieldLevel) bool {
	Image, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}
	// format Check
	var ImageFormats = []string{"png", "jpg", "jpeg", "webp"}
	var formatCheck bool
	for _, item := range ImageFormats {
		if Image.Filename[len(Image.Filename)-3:] == item {
			formatCheck = true
		}
	}
	if !formatCheck {
		return false
	}
	// Size > 6mb
	if Image.Size > 6000 {
		return false
	}
	return true
}
