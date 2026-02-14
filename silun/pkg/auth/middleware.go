package auth

import (
	"context"
	"silun/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

func AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.GetHeader("Access-Token")
		if token == "" {
			utils.Error(c, -1, "Unauthorized")
			c.Abort()
			return
		}

		userID, err := ParseAccessToken(token)
		if err != nil {
			utils.Error(c, -1, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next(ctx)
	}
}

func ParseAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-access-secret-key"), nil
	}, &Claims{})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
