package apperror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleError(c *gin.Context, logger *zap.Logger, err error) {
	var appErr *Error
	if !errors.As(err, &appErr) {
		appErr = Internal(err)
	}

	status := http.StatusInternalServerError
	switch appErr.Code {
	case CodeValidation:
		status = http.StatusBadRequest
	case CodeNotFound:
		status = http.StatusNotFound
	case CodeConflict:
		status = http.StatusConflict
	}

	if appErr.Err != nil {
		logger.Error(appErr.Message, zap.Error(appErr.Err))
	}
	c.AbortWithStatusJSON(status, gin.H{"error": appErr.Message})
}
