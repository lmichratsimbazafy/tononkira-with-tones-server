package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthTokenPayload struct {
	access_token  string
	refresh_token string
}

type AuthTokenClaims struct {
	jwt.StandardClaims
	Sub        primitive.ObjectID `json:"sub"`
	UserName   string             `json:"userName"`
	Permission interface{}        `json:"permission"`
}

func getToken(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	return token[len("Bearer "):]
}

// func VerifyToken(token string) (interface{}, error) {
// 	claims := AuthTokenClaims{}
// 	_, err := jwt.ParseWithClaims(token, &claims, func(_token *jwt.Token) (interface{}, error) {
// 		_, ok := _token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", _token.Header["alg"])
// 		}
// 		return []byte(config.Env.JWTSecret), nil
// 	})
// 	if err != nil {
// 		authToken := models
// 		if err == jwt.ErrECDSAVerification {}
// 		return nil, ErrorWrapper(err)
// 	}
// 	// fmt.Printf("claims %v", claims.Permission.Lyrics)
// 	// if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
// 	// 	// Get the user record from database or
// 	// 	// run through your business logic to verify if the user can log in
// 	// 	if int(claims["sub"].(float64)) == 1 {

// 	// 		return c.JSON(http.StatusOK, newTokenPair)
// 	// 	}

//		// 	return echo.ErrUnauthorized
//		// }
//		return claims, nil
//	}

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
