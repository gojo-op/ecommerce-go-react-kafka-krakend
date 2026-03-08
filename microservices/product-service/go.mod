module product-service

go 1.21

require (
    github.com/gin-gonic/gin v1.11.0
    github.com/google/uuid v1.6.0
    github.com/glebarez/sqlite v1.11.0
    gorm.io/gorm v1.25.7
    github.com/IBM/sarama v1.43.3
)

replace github.com/glebarez/sqlite => github.com/glebarez/sqlite v1.11.0