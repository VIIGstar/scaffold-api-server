package user

import "scaffold-api-server/internal/entities"

func (i impl) Create(u *entities.Investor) error {
	return i.db.GormDB().Model(entities.Investor{}).Create(u).Error
}
