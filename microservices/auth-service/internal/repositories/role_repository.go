package repositories

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "auth-service/internal/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) FindByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Role{}, "id = ?", id).Error
}

func (r *RoleRepository) FindAll(limit, offset int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.Preload("Permissions")
	
	if err := query.Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepository) AssignPermission(roleID, permissionID uuid.UUID) error {
	return r.db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error
}

func (r *RoleRepository) RevokePermission(roleID, permissionID uuid.UUID) error {
	return r.db.Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID).Error
}

func (r *RoleRepository) GetRolePermissions(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Select("permissions.*").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

func (r *RoleRepository) FindByUserID(userID uuid.UUID) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}