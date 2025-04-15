package utils

import (
	"fmt"
	"mime"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
)

func AllowedTypeByExtension(
	extension string,
	allowedAvatarContentTypes map[string]string,
) (string, error) {
	contentType := mime.TypeByExtension(extension)

	_, ok := allowedAvatarContentTypes[extension]
	if !ok {
		return contentType, apperrs.NewError(
			fmt.Errorf("unknown content type for extension, %v", extension),
		)
	}

	return contentType, nil
}

func GetContentTypeFromExt(
	extension string,
) string {
	return mime.TypeByExtension(extension)
}
