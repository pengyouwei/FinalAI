package myvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	// 使用 json tag 作为字段名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	return &CustomValidator{validator: v}
}

// 核心：实现 Echo Validator 接口
func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {
			var msg strings.Builder

			for _, e := range errs {
				field := e.Field() // 已经是 json tag
				tag := e.Tag()     // 校验规则
				param := e.Param() // 参数（如 min=3 → 3）

				msg.WriteString(formatError(field, tag, param))
				msg.WriteString("; ")
			}

			return fmt.Errorf("%s", msg.String())
		}

		return err
	}

	return nil
}

func formatError(field, tag, param string) string {
	switch tag {

	case "required":
		return fmt.Sprintf("%s 为必填项", field)

	case "min":
		return fmt.Sprintf("%s 最小长度为 %s", field, param)

	case "max":
		return fmt.Sprintf("%s 最大长度为 %s", field, param)

	case "email":
		return fmt.Sprintf("%s 格式不正确", field)

	case "oneof":
		return fmt.Sprintf("%s 必须是 [%s] 之一", field, param)

	case "len":
		return fmt.Sprintf("%s 长度必须为 %s", field, param)

	case "numeric":
		return fmt.Sprintf("%s 必须为数字", field)

	case "eqfield":
		return fmt.Sprintf("%s 必须与 %s 相同", field, param)

	default:
		return fmt.Sprintf("%s 不合法(%s)", field, tag)
	}
}
