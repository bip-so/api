package mentions

import "gitlab.com/phonepost/bip-be-platform/internal/models"

func (r mentionsRepo) GetRoleObjects(rolesIDs []uint64) (*[]models.Role, error) {
	var roles *[]models.Role
	err := r.db.Model(&models.Role{}).Where("id IN ?", rolesIDs).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r mentionsRepo) GetBranchObjects(branchesIDs []uint64) (*[]models.CanvasBranch, error) {
	var branches *[]models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("id IN ?", branchesIDs).Find(&branches).Error
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (r mentionsRepo) GetBranchObject(branchID uint64) (*models.CanvasBranch, error) {
	var branch *models.CanvasBranch
	err := r.db.Model(&models.CanvasBranch{}).Where("id = ?", branchID).First(&branch).Error
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (r mentionsRepo) GetUserObjects(usersIDs []uint64) (*[]models.User, error) {
	var users *[]models.User
	err := r.db.Model(&models.User{}).Where("id IN ?", usersIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
