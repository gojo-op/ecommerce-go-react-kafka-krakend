package middleware

import (
  "fmt"
  "net/http"
  "strings"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct{
  UserID string `json:"user_id"`
  Email string `json:"email"`
  Username string `json:"username"`
  Roles []string `json:"roles"`
  Permissions []string `json:"permissions"`
  jwt.RegisteredClaims
}

type AuthMiddleware struct{ jwtSecret []byte }
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware { return &AuthMiddleware{ jwtSecret: []byte(jwtSecret) } }

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
  return func(c *gin.Context){
    auth := c.GetHeader("Authorization")
    if auth == "" { c.JSON(http.StatusUnauthorized, gin.H{"error":"Authorization header required"}); c.Abort(); return }
    parts := strings.Split(auth, " ")
    if len(parts) != 2 || parts[0] != "Bearer" { c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid authorization header format"}); c.Abort(); return }
    claims, err := a.validateToken(parts[1])
    if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid or expired token"}); c.Abort(); return }
    c.Set("user_id", claims.UserID)
    c.Set("email", claims.Email)
    c.Set("username", claims.Username)
    c.Set("roles", claims.Roles)
    c.Set("permissions", claims.Permissions)
    c.Set("claims", claims)
    c.Next()
  }
}

func (a *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
  return func(c *gin.Context){ a.RequireAuth()(c); if c.IsAborted(){ return }
    v, ok := c.Get("roles"); if !ok { c.JSON(http.StatusForbidden, gin.H{"error":"No roles found"}); c.Abort(); return }
    roles, _ := v.([]string)
    has := false
    for _, r := range roles { if r == role { has = true; break } }
    if !has { c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Role '%s' required", role)}); c.Abort(); return }
    c.Next()
  }
}

func (a *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
  token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]) }
    return a.jwtSecret, nil
  })
  if err != nil { return nil, err }
  claims, ok := token.Claims.(*JWTClaims)
  if !ok || !token.Valid { return nil, fmt.Errorf("invalid token claims") }
  if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) { return nil, fmt.Errorf("token has expired") }
  return claims, nil
}