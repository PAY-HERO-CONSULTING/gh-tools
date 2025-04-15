package sessions

import (
	"github.com/gin-gonic/gin"
)

type (
	RequestValidator interface {
		ValidateOrigin(c *gin.Context) (string, error)
	}
)
