package ex

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		log.Printf("[%d] - %s in %v\n", ctx.StatusCode, ctx.Path, time.Since(start))
	}
}

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recovered: ", err)
				ctx.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

func CORS(config CORSConfig) HandlerFunc {
	return func(ctx *Context) {
		origin := ctx.Req.Header.Get("Origin")
		allowed := false
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
				ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		if len(config.AllowMethods) > 0 {
			methods := ""
			for i, m := range config.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += m
			}
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		}

		if len(config.AllowHeaders) > 0 {
			headers := ""
			for i, h := range config.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += h
			}
			ctx.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		}

		if len(config.ExposeHeaders) > 0 {
			headers := ""
			for i, h := range config.ExposeHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += h
			}
			ctx.Writer.Header().Set("Access-Control-Expose-Headers", headers)
		}

		if config.AllowCredentials {
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if config.MaxAge > 0 {
			ctx.Writer.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		if ctx.Method == "OPTIONS" {
			ctx.Status(http.StatusNoContent)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func DefaultCORS() HandlerFunc {
	return CORS(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	})
}

const RequestIDHeader = "X-Request-ID"

func RequestID() HandlerFunc {
	return func(ctx *Context) {
		requestID := ctx.Req.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}
		ctx.Writer.Header().Set(RequestIDHeader, requestID)
		ctx.Next()
	}
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(bytes)
}

type JWTConfig struct {
	Secret      string
	TokenLookup string
	AuthScheme  string
	ContextKey  string
}

type JWTClaims struct {
	Issuer    string                 `json:"iss"`
	Subject   string                 `json:"sub"`
	Audience  string                 `json:"aud"`
	ExpiresAt int64                  `json:"exp"`
	NotBefore int64                  `json:"nbf"`
	IssuedAt  int64                  `json:"iat"`
	ID        string                 `json:"jti"`
	Payload   map[string]interface{} `json:"payload"`
}

type jwtToken struct {
	Header    map[string]interface{} `json:"header"`
	Payload   JWTClaims              `json:"payload"`
	Signature string                 `json:"signature"`
}

func JWT(config JWTConfig) HandlerFunc {
	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization"
	}
	if config.AuthScheme == "" {
		config.AuthScheme = "Bearer"
	}
	if config.ContextKey == "" {
		config.ContextKey = "user"
	}

	return func(ctx *Context) {
		tokenStr := extractToken(ctx, config.TokenLookup, config.AuthScheme)
		if tokenStr == "" {
			ctx.Status(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "missing or malformed jwt")
			ctx.Abort()
			return
		}

		claims, err := parseToken(tokenStr, config.Secret)
		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}

		if claims.ExpiresAt > 0 && time.Now().Unix() > claims.ExpiresAt {
			ctx.Status(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "token expired")
			ctx.Abort()
			return
		}

		if claims.NotBefore > 0 && time.Now().Unix() < claims.NotBefore {
			ctx.Status(http.StatusUnauthorized)
			ctx.String(http.StatusUnauthorized, "token not valid yet")
			ctx.Abort()
			return
		}

		ctx.Req = ctx.Req.WithContext(setValue(ctx.Req.Context(), config.ContextKey, claims))
		ctx.Next()
	}
}

func extractToken(ctx *Context, lookup, authScheme string) string {
	parts := strings.Split(lookup, ":")
	if len(parts) != 2 {
		return ""
	}

	switch parts[0] {
	case "header":
		auth := ctx.Req.Header.Get(parts[1])
		if auth == "" {
			return ""
		}
		if authScheme != "" {
			if strings.HasPrefix(auth, authScheme+" ") {
				return strings.TrimPrefix(auth, authScheme+" ")
			}
			return ""
		}
		return auth
	case "query":
		return ctx.Query(parts[1])
	default:
		return ""
	}
}

func parseToken(tokenStr, secret string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, &jwtError{message: "invalid token format"}
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, &jwtError{message: "invalid header encoding"}
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, &jwtError{message: "invalid header format"}
	}

	alg, ok := header["alg"].(string)
	if !ok || alg != "HS256" {
		return nil, &jwtError{message: "unsupported algorithm"}
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, &jwtError{message: "invalid payload encoding"}
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, &jwtError{message: "invalid payload format"}
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := signHS256(signingInput, secret)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, &jwtError{message: "invalid signature"}
	}

	return &claims, nil
}

func signHS256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func GenerateToken(claims JWTClaims, secret string) (string, error) {
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signingInput := headerEncoded + "." + payloadEncoded
	signature := signHS256(signingInput, secret)

	return signingInput + "." + signature, nil
}

type jwtError struct {
	message string
}

func (e *jwtError) Error() string {
	return e.message
}

type contextKey string

func setValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

func GetJWTClaims(ctx *Context, key string) *JWTClaims {
	if key == "" {
		key = "user"
	}
	if claims, ok := ctx.Req.Context().Value(contextKey(key)).(*JWTClaims); ok {
		return claims
	}
	return nil
}
