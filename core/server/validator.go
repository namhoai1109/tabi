package server

import (
	"io"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/go-playground/validator/v10"
)

// custom variable
var (
	DocumentExts = []string{".pdf", ".jpg", ".jpeg", ".png"}
	ImageExts    = []string{".jpg", ".jpeg", ".png"}
)

// CustomValidator holds custom validator
type CustomValidator struct {
	V *validator.Validate
}

// NewValidator creates new custom validator
func NewValidator() *CustomValidator {
	V := validator.New()
	V.RegisterValidation("date", validateDate)
	V.RegisterValidation("mobile", validateMobile)
	V.RegisterValidation("document", validateDocument)
	V.RegisterValidation("image", validateImage)
	V.RegisterValidation("fullname", validateFullname)
	V.RegisterValidation("description", validateDescription)
	return &CustomValidator{V}
}

// Validate validates the request
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.V.Struct(i)
}

func validateDate(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	re := regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}(T00:00:00Z)?$`)
	return re.MatchString(val)
}

func validateMobile(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	re := regexp.MustCompile(`^([0]?(3|5|7|8|9|1[2|6|8|9]))([0-9]{8})\b`)
	return re.MatchString(strings.Replace(val, " ", "", -1))
}

func validateFullname(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	re := regexp.MustCompile(`^([a-zA-Z\xC0-\uFFFF]{1,40}[ \-\‘]{0,1}){1,10}$`)
	return re.MatchString(strings.Replace(val, " ", "", -1))
}

func validateDescription(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	re := regexp.MustCompile(`^([a-zA-Z0-9\xC0-\uFFFF]{1,40}[ \-\’]{0,1}){1,10}$`)
	return re.MatchString(strings.Replace(val, " ", "", -1))
}

func validateDocument(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().([]*multipart.FileHeader)
	if len(val) == 0 {
		return true
	}
	f := val[0]

	src, err := f.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// read file
	buff, err := io.ReadAll(src)
	if err != nil {
		return false
	}
	t := mimetype.Detect(buff)
	for _, v := range DocumentExts {
		if v == t.Extension() {
			return true
		}
	}
	return false
}

func validateImage(fl validator.FieldLevel) bool {
	val := fl.Field().Interface().([]*multipart.FileHeader)
	if len(val) == 0 {
		return true
	}
	f := val[0]
	src, err := f.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	// read file
	buff, err := io.ReadAll(src)
	if err != nil {
		return false
	}
	t := mimetype.Detect(buff)
	for _, v := range ImageExts {
		if v == t.Extension() {
			return true
		}
	}
	return false
}
