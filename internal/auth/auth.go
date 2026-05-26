package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/lookingcamel/system-framework/internal/config"
	"github.com/lookingcamel/system-framework/internal/logger"
	"github.com/lookingcamel/system-framework/pkg/utils"
)

var cfg config.AuthConfig

func Init(config config.AuthConfig) {
	cfg = config
	if cfg.Enabled {
		logger.Log.Info("Authentication initialized",
			zap.Bool("jwt_enabled", cfg.JWTEnabled),
			zap.Bool("signature_enabled", cfg.SignatureEnabled),
		)
	}
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, username, role string) (string, error) {
	expireTime := time.Now().Add(time.Second * time.Duration(cfg.JWTExpireSeconds))

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWTIssuer,
			Subject:   userID,
			Audience:  []string{cfg.JWTAudience},
			ExpiresAt: jwt.NewNumericDate(expireTime),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        utils.GenerateRequestID(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled || !cfg.JWTEnabled {
			c.Next()
			return
		}

		if isExcludedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := ParseToken(tokenString)
		if err != nil {
			logger.Log.Warn("Invalid token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Invalid token: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		ctx := context.WithValue(c.Request.Context(), "claims", claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func Signature() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled || !cfg.SignatureEnabled {
			c.Next()
			return
		}

		if isExcludedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		signature := c.GetHeader(cfg.SignatureHeader)
		timestamp := c.GetHeader(cfg.TimestampHeader)
		nonce := c.GetHeader(cfg.NonceHeader)

		if signature == "" || timestamp == "" {
			logger.Log.Warn("Missing signature or timestamp")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Missing signature or timestamp",
			})
			c.Abort()
			return
		}

		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			logger.Log.Warn("Invalid timestamp format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Invalid timestamp format",
			})
			c.Abort()
			return
		}

		now := time.Now().Unix()
		if abs(now-timestampInt) > int64(cfg.TimestampMaxDiff) {
			logger.Log.Warn("Timestamp expired",
				zap.Int64("now", now),
				zap.Int64("timestamp", timestampInt),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Request expired",
			})
			c.Abort()
			return
		}

		if err := validateNonce(nonce); err != nil {
			logger.Log.Warn("Invalid nonce", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		expectedSignature := generateSignature(c, timestamp, nonce)
		if signature != expectedSignature {
			logger.Log.Warn("Signature mismatch")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Invalid signature",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateSignature(c *gin.Context, timestamp, nonce string) string {
	body, _ := c.GetRawData()
	c.Request.Body = nil

	data := fmt.Sprintf("%s%s%s%s",
		c.Request.Method,
		c.Request.URL.Path,
		timestamp,
		string(body),
	)

	if nonce != "" {
		data += nonce
	}

	mac := hmac.New(sha256.New, []byte(cfg.SignatureSecret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func validateNonce(nonce string) error {
	if nonce == "" {
		return nil
	}

	if len(nonce) < 16 {
		return fmt.Errorf("nonce too short")
	}

	return nil
}

func isExcludedPath(path string) bool {
	for _, excludePath := range cfg.ExcludePaths {
		if strings.HasPrefix(path, excludePath) {
			return true
		}
	}
	return false
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

func GetRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}

func GetClaims(c *gin.Context) *Claims {
	if claims, exists := c.Get("claims"); exists {
		return claims.(*Claims)
	}
	return nil
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetRole(c)
		if userRole != role {
			logger.Log.Warn("Insufficient permissions",
				zap.String("required_role", role),
				zap.String("user_role", userRole),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    -1,
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetRole(c)
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		logger.Log.Warn("Insufficient permissions",
			zap.String("user_role", userRole),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"code":    -1,
			"message": "Insufficient permissions",
		})
		c.Abort()
	}
}