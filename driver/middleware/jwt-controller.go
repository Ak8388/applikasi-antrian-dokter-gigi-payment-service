package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/utils/common"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	JwtVerify(role ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtUtils common.JwtToken
}

func (am *authMiddleware) JwtVerify(role ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.Request.Header.Get("Authorization")
		tokenString := strings.Replace(authorization, "Bearer ", "", -1)

		payloadToken := model.TokenAkses{
			TokenString: tokenString,
		}

		claims, err := am.jwtUtils.VerfifyToken(payloadToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		validateRole := false

		for _, x := range role {
			if x == claims["Role"].(string) {
				validateRole = true
			}
		}

		if !validateRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "can't akses this page. invalid role"})
			return
		}

		exp := claims["exp"].(float64)

		if time.Now().After(time.Unix(int64(exp), 10)) {
			ctx.AbortWithStatusJSON(419, gin.H{"error": "Hey, your session has ended"})
			return
		}

		ctx.Set("userEmail", claims["Email"].(string))
		ctx.Set("userRole", claims["Role"].(string))
		ctx.Set("userID", claims["Id"].(string))
		ctx.Next()
	}
}

func NewAuthMiddleware(jwtUtils common.JwtToken) AuthMiddleware {
	return &authMiddleware{jwtUtils}
}
