package utils

import (
	"context"
	"errors"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

func GetErrorStatus(
	appError interface{},
) int {
	if newAppErr, ok := appError.(*apperrs.Error); ok {
		return newAppErr.Status()
	} else {
		defaultAppErr := apperrs.New(errors.New("unknown error"), apperrs.UnexpextedError)
		return defaultAppErr.Status()
	}
}

func HandleAPIError(
	c *gin.Context,
	appError *apperrs.Error,
	service string,
) {
	logger.Errorf(service, "failed to make request err: [%v]", appError)
	c.JSON(appError.Status(), appError.JsonResponse())
}

func HandleError(
	c *gin.Context,
	appError *apperrs.Error,
) {
	logger.AppError(c.Request.Context(), appError)
	c.JSON(appError.Status(), appError.JsonResponse())
}

func RPCError(
	ctx context.Context,
	err error,
) error {
	appErr := apperrs.NewError(err)

	logger.AppError(ctx, appErr)

	return status.Errorf(appErr.GRPCStatus(), appErr.GetMessage())
}
