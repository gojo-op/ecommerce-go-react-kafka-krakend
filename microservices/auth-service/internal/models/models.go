package models

import (
  "time"
  "github.com/google/uuid"
  "gorm.io/gorm"
)

type User struct {
  ID              uuid.UUID      `json:"id" gorm:"type:text;primaryKey"`
  Email           string         `json:"email" gorm:"uniqueIndex;not null"`
  Username        string         `json:"username" gorm:"uniqueIndex;not null"`
  Password        string         `json:"-" gorm:"not null"`
  FirstName       string         `json:"first_name"`
  LastName        string         `json:"last_name"`
  Phone           string         `json:"phone"`
  AvatarURL       string         `json:"avatar_url"`
  IsEmailVerified bool           `json:"is_email_verified" gorm:"default:false"`
  IsPhoneVerified bool           `json:"is_phone_verified" gorm:"default:false"`
  IsActive        bool           `json:"is_active" gorm:"default:true"`
  LastLoginAt     *time.Time     `json:"last_login_at"`
  CreatedAt       time.Time      `json:"created_at"`
  UpdatedAt       time.Time      `json:"updated_at"`
  DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

  Roles       []Role       `json:"roles,omitempty" gorm:"many2many:user_roles;"`
  Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:user_permissions;"`
  Profile     *UserProfile `json:"profile,omitempty"`
  Addresses   []Address    `json:"addresses,omitempty"`
}

type UserProfile struct {
  ID          uuid.UUID  `json:"id" gorm:"type:text;primaryKey"`
  UserID      uuid.UUID  `json:"user_id" gorm:"type:text;not null;uniqueIndex"`
  Bio         string     `json:"bio" gorm:"type:text"`
  DateOfBirth *time.Time `json:"date_of_birth"`
  Gender      string     `json:"gender" gorm:"type:varchar(10)"`
  Country     string     `json:"country" gorm:"type:varchar(50)"`
  City        string     `json:"city" gorm:"type:varchar(50)"`
  Timezone    string     `json:"timezone" gorm:"type:varchar(50)"`
  Language    string     `json:"language" gorm:"type:varchar(10);default:'en'"`
  CreatedAt   time.Time  `json:"created_at"`
  UpdatedAt   time.Time  `json:"updated_at"`

  User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Role struct {
  ID          uuid.UUID      `json:"id" gorm:"type:text;primaryKey"`
  Name        string         `json:"name" gorm:"uniqueIndex;not null"`
  Description string         `json:"description"`
  IsActive    bool           `json:"is_active" gorm:"default:true"`
  CreatedAt   time.Time      `json:"created_at"`
  UpdatedAt   time.Time      `json:"updated_at"`
  DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

  Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
  Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
}

type Permission struct {
  ID          uuid.UUID      `json:"id" gorm:"type:text;primaryKey"`
  Name        string         `json:"name" gorm:"uniqueIndex;not null"`
  Resource    string         `json:"resource" gorm:"not null"`
  Action      string         `json:"action" gorm:"not null"`
  Description string         `json:"description"`
  CreatedAt   time.Time      `json:"created_at"`
  UpdatedAt   time.Time      `json:"updated_at"`
  DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

  Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
  Users []User `json:"users,omitempty" gorm:"many2many:user_permissions;"`
}

type Address struct {
  ID         uuid.UUID `json:"id" gorm:"type:text;primaryKey"`
  UserID     uuid.UUID `json:"user_id" gorm:"type:text;not null;index"`
  Type       string    `json:"type" gorm:"type:varchar(20);not null"`
  FirstName  string    `json:"first_name" gorm:"not null"`
  LastName   string    `json:"last_name" gorm:"not null"`
  Company    string    `json:"company"`
  Address1   string    `json:"address1" gorm:"not null"`
  Address2   string    `json:"address2"`
  City       string    `json:"city" gorm:"not null"`
  State      string    `json:"state" gorm:"not null"`
  Country    string    `json:"country" gorm:"not null"`
  PostalCode string    `json:"postal_code" gorm:"not null"`
  Phone      string    `json:"phone"`
  IsDefault  bool      `json:"is_default" gorm:"default:false"`
  CreatedAt  time.Time `json:"created_at"`
  UpdatedAt  time.Time `json:"updated_at"`

  User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type TokenPair struct {
  AccessToken  string `json:"access_token"`
  RefreshToken string `json:"refresh_token"`
  ExpiresAt    int64  `json:"expires_at"`
}

type LoginRequest struct {
  Email    string `json:"email" binding:"required,email"`
  Password string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
  Email     string `json:"email" binding:"required,email"`
  Username  string `json:"username" binding:"required,min=3,max=20"`
  Password  string `json:"password" binding:"required,min=8"`
  FirstName string `json:"first_name" binding:"required"`
  LastName  string `json:"last_name" binding:"required"`
  Phone     string `json:"phone"`
}

type UpdateProfileRequest struct {
  FirstName   string `json:"first_name"`
  LastName    string `json:"last_name"`
  Phone       string `json:"phone"`
  Bio         string `json:"bio"`
  DateOfBirth string `json:"date_of_birth"`
  Gender      string `json:"gender"`
  Country     string `json:"country"`
  City        string `json:"city"`
  Timezone    string `json:"timezone"`
  Language    string `json:"language"`
}

type ChangePasswordRequest struct {
  OldPassword string `json:"old_password" binding:"required"`
  NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserResponse struct {
  ID              uuid.UUID      `json:"id"`
  Email           string         `json:"email"`
  Username        string         `json:"username"`
  FirstName       string         `json:"first_name"`
  LastName        string         `json:"last_name"`
  Phone           string         `json:"phone"`
  AvatarURL       string         `json:"avatar_url"`
  IsEmailVerified bool           `json:"is_email_verified"`
  IsPhoneVerified bool           `json:"is_phone_verified"`
  IsActive        bool           `json:"is_active"`
  LastLoginAt     *time.Time     `json:"last_login_at"`
  CreatedAt       time.Time      `json:"created_at"`
  UpdatedAt       time.Time      `json:"updated_at"`
  Roles           []RoleResponse `json:"roles,omitempty"`
  Profile         *UserProfile   `json:"profile,omitempty"`
}

type RoleResponse struct {
  ID          string `json:"id"`
  Name        string `json:"name"`
  Description string `json:"description"`
}

const (
  RoleAdmin     = "admin"
  RoleUser      = "user"
  RoleSeller    = "seller"
  RoleCustomer  = "customer"
  RoleModerator = "moderator"
  RoleSupport   = "support"
)