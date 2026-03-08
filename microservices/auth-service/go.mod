module auth-service

go 1.21

require (
    github.com/gin-gonic/gin v1.11.0
    github.com/golang-jwt/jwt/v5 v5.3.0
    github.com/google/uuid v1.6.0
    gorm.io/gorm v1.25.7
    github.com/glebarez/sqlite v1.11.0
)

replace github.com/glebarez/sqlite => github.com/glebarez/sqlite v1.11.0
