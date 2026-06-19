package middleware

import (
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	CtxKeyUserID = "user_id"
	CtxKeyUser   = "current_user"
	CookieName   = "auth_token"
)

func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr, err := ctx.Cookie(CookieName)
		if err != nil || tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth cookie"})
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invlid token"})
			return
		}

		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			return
		}

		ctx.Set(CtxKeyUserID, userID)
		ctx.Next()
	}
}

func LoadCurrentUser(userUC *usecase.UserUsecase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIDVal, exists := ctx.Get(CtxKeyUserID)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		userID := userIDVal.(uuid.UUID)

		user, err := userUC.FindByID(ctx.Request.Context(), userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		ctx.Set(CtxKeyUser, user)
		ctx.Next()
	}
}

func RequireRole(role entity.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := ctx.Get(CtxKeyUser)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		if user.(*entity.User).Role != role {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		ctx.Next()
	}
}

func CurrentUserID(ctx *gin.Context) uuid.UUID {
	return ctx.MustGet(CtxKeyUserID).(uuid.UUID)
}

func CurrentUser(ctx *gin.Context) *entity.User {
	return ctx.MustGet(CtxKeyUser).(*entity.User)
}
