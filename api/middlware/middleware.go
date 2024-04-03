package middlware

import (
	"testing-api-gateway/api/handlers/tokens"

	"github.com/gin-gonic/gin"

	// "github.com/spf13/cast"

	"net/http"
)

func Auth(ctx *gin.Context) {

	if ctx.Request.URL.Path == "/v1/verification" || ctx.Request.URL.Path == "/v1/login" || ctx.Request.URL.Path == "/v1/register" || ctx.Request.URL.Path == "/v1/swagger/swaggerdoc.json" || ctx.Request.URL.Path == "/v1/swagger/index.html" || ctx.Request.URL.Path == "/v1/swagger/swagger-ui.css" || ctx.Request.URL.Path == "/v1/swagger/swagger-ui-bundle.js" || ctx.Request.URL.Path == "/v1/swagger/swagger-ui-standalone-preset.js" {
		ctx.Next()
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized request",
		})
		return
	}

	_, err := tokens.ExtractClaim(token, []byte("key"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
}
