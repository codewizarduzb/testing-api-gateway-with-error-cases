package casbin

import (
	"errors"
	"net/http"
	"strings"
	"testing-api-gateway/api/handlers/tokens"
	"testing-api-gateway/config"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type CasbinHandler struct {
	config   config.Config
	enforcer *casbin.Enforcer
	jwt      tokens.JWTHandler
}

func NewAuthorizer() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.Next()
			return
		}
		claims, err := tokens.ExtractClaim(token, []byte("key"))

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}
		sub := claims["role"]
		obj := ctx.Request.URL.Path
		etc := ctx.Request.Method
		e, _ := casbin.NewEnforcer("auth.conf", "auth.csv")
		t, _ := e.Enforce(sub, obj, etc)
		if t {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "no access ",
		})
	}
}

func CheckCasbinPermission(casbin *casbin.Enforcer, config config.Config) gin.HandlerFunc {
	casbHandler := &CasbinHandler{
		config:   config,
		enforcer: casbin,
	}
	return func(ctx *gin.Context) {
		allowed, err := casbHandler.CheckPermission(ctx.Request)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
	}
}

func (casbin *CasbinHandler) GetRole(ctx *http.Request) (string, int) {
	var r string
	token := ctx.Header.Get("Authorization")
	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		r = strings.TrimPrefix(token, "Bearer ")
	} else {
		r = token
	}
	claims, err := tokens.ExtractClaim(r, []byte("key"))
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["role"]), 0
}

func (casbin *CasbinHandler) CheckPermission(r *http.Request) (bool, error) {
	role, status := casbin.GetRole(r)
	if role == "unauthorized" {
		return true, nil
	}
	if status != 0 {
		return false, errors.New(role)
	}
	method := r.Method
	action := r.URL.Path
	c, err := casbin.enforcer.Enforce(role, action, method)
	if err != nil {
		return false, err
	}

	return c, nil
}
