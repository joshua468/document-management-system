package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joshua468/document-management-system/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoToken          = errors.New("no token provided")
	ErrInvalidToken     = errors.New("invalid token")
	ErrUnauthorizedUser = errors.New("unauthorized user")
	userIDContextKey    = "userID"
)

type Service struct {
	userRepo user.UserRepository
}

func NewService(userRepo user.UserRepository) *Service {
	return &Service{userRepo}
}

func (s *Service) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrUnauthorizedUser
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("secret"))
}

func (s *Service) ValidateToken(r *http.Request) (uint, error) {
	tokenString := ExtractTokenFromHeader(r)
	if tokenString == "" {
		return 0, ErrNoToken
	}
	return s.VerifyToken(tokenString)
}

func ExtractTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token == "" {
		return ""
	}
	return token[7:]
}

func (s *Service) VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}
	userID := uint(claims["userID"].(float64))
	return userID, nil
}

func SetUserIDInContext(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetUserIDFromContext(ctx context.Context) uint {
	userID, _ := ctx.Value(userIDContextKey).(uint)
	return userID
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// HandleLogin handles the login request
func HandleLogin(authService *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		token, err := authService.Login(credentials.Username, credentials.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(TokenResponse{Token: token})
	}
}
