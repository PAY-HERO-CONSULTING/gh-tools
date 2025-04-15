package utils

import (
	"reflect"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArguement struct {
	Field string `json:"field"`
	Value any    `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func Bind(c *gin.Context, req any) bool {
	if err := c.ShouldBind(req); err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArguement

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArguement{
					getTagName(reflect.TypeOf(req).Elem(), err.Field()),
					err.Value(),
					err.Tag(),
					err.Param(),
				})
			}

			appErr := apperrs.NewBadRequest("Submitted invalid form")

			response := appErr.JsonResponse()
			response["invalid_arguments"] = invalidArgs

			c.JSON(appErr.Status(), response)

			return false
		}

		appErr := apperrs.Wrap(err)
		logger.AppError(c.Request.Context(), appErr)

		c.JSON(appErr.Status(), appErr.JsonResponse())

		return false
	}

	return true
}

func getTagName(t reflect.Type, fieldName string) string {
	field, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}

	tag := field.Tag.Get("json")
	if tag != "" {
		return tag
	}

	tag = field.Tag.Get("form")
	if tag != "" {
		return tag
	}

	return field.Tag.Get("xml")
}
