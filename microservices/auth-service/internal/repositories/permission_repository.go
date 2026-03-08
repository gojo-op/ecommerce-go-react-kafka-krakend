package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "auth-service/internal/models"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) FindByID(id uuid.UUID) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) FindByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Permission{}, "id = ?", id).Error
}

func (r *PermissionRepository) FindAll(limit, offset int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	if err := r.db.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

func (r *PermissionRepository) FindByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("resource = ?", resource).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindByRoleID(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Select("permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindByUserID(userID uuid.UUID) ([]models.Permission, error) {
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