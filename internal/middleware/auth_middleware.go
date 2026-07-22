package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/erenmentes/go-prod-ready-auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type authMiddleware struct {
	db *gorm.DB
}

type jwtClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type contextKey string

const userContextKey = contextKey("user")

func NewAuthMiddleware(db *gorm.DB) *authMiddleware {
	return &authMiddleware{db: db}
}

func (m *authMiddleware) Apply(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			respondError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
		claims := &jwtClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		var user models.User
		err = m.db.Where("id = ?", claims.UserID).First(&user).Error
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid session")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userContextKey, &user))
		next(w, r)
	}
}

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	value := ctx.Value(userContextKey)
	user, ok := value.(*models.User)
	return user, ok
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
