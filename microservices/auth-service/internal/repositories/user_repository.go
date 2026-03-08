package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "auth-service/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uuid.UUID) (models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").Preload("Permissions").Preload("Profile").First(&user, "id = ?", id).Error
	return user, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").Preload("Permissions").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").Preload("Permissions").First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateProfile(profile *models.UserProfile) error {
	return r.db.Save(profile).Error
}

func (r *UserRepository) CreateProfile(profile *models.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *UserRepository) FindAll(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Preload("Roles.Permissions").Preload("Permissions").Preload("Profile")
	
	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FindByRole(roleID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Roles.Permissions").Preload("Permissions").Preload("Profile").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID).Find(&users).Error
	return users, err
}

func (r *UserRepository) AssignRole(userID, roleID uuid.UUID) error {
	return r.db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
}

func (r *UserRepository) RevokeRole(userID, roleID uuid.UUID) error {
	return r.db.Exec("DELETE FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Error
}

func (r *UserRepository) AssignPermission(userID, permissionID uuid.UUID) error {
	return r.db.Exec("INSERT INTO user_permissions (user_id, permission_id) VALUES (?, ?)", userID, permissionID).Error
}

func (r *UserRepository) RevokePermission(userID, permissionID uuid.UUID) error {
	return r.db.Exec("DELETE FROM user_permissions WHERE user_id = ? AND permission_id = ?", userID, permissionID).Error
}

func (r *UserRepository) GetUserPermissions(userID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	
	// Get permissions from roles
	err := r.db.Table("permissions").
		Select("DISTINCT permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&permissions).Error
	
	if err != nil {
		return nil, err
	}

	// Get direct permissions
	var directPermissions []models.Permission
	err = r.db.Table("permissions").
		Select("DISTINCT permissions.*").
		Joins("JOIN user_permissions ON permissions.id = user_permissions.permission_id").
		Where("user_permissions.user_id = ?", userID).
		Find(&directPermissions).Error
	
	if err != nil {
		return nil, err
	}

	// Combine and remove duplicates
	permissionMap := make(map[uuid.UUID]models.Permission)
	for _, perm := range permissions {
		permissionMap[perm.ID] = perm
	}
	for _, perm := range directPermissions {
		permissionMap[perm.ID] = perm
	}

	result := make([]models.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		result = append(result, perm)
	}

	return result, nil
}