package sessions

import (
	"github.com/PAY-HERO-CONSULTING/gh-tools/models"
	"github.com/gin-gonic/gin"
)

type (
	AuthSession interface {
		TokenInfo(c *gin.Context) (*models.TokenInfo, error)
	}
)
