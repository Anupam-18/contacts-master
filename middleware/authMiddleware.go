package middleware

import (
	"contact-store/models"
	"contact-store/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var HandleJwtAuth = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("Authorization")
		if len(strings.TrimSpace(headerToken)) == 0 {
			utils.Respond(w, 401, utils.Message("401", "Missing authorization"))
			return
		}
		splittedToken := strings.Split(headerToken, " ")
		if len(splittedToken) != 2 {
			fmt.Println("came here")
			utils.Respond(w, 401, utils.Message("401", "Malformed authorization token"))
			return
		}
		tokenObj := &models.Token{}
		token, err := jwt.ParseWithClaims(splittedToken[1], tokenObj, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			fmt.Println("came here", err)
			utils.Respond(w, 403, utils.Message("403", "Malformed authorization token"))
			return
		}
		if !token.Valid {
			utils.Respond(w, 403, utils.Message("403", "Invalid token"))
			return
		}
		ctx := context.WithValue(r.Context(), "user", tokenObj.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
