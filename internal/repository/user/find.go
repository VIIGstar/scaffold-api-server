package user

import (
	"context"
	"gorm.io/gorm"
	"scaffold-api-server/internal/entities"
	query_params "scaffold-api-server/internal/query-params"
	"scaffold-api-server/pkg/database"
	info_log "scaffold-api-server/pkg/info-log"
)

func (i impl) Find(ctx context.Context, req query_params.GetUserParams, lock bool) (*entities.Investor, error) {
	if !isValidParams(req) {
		return nil, database.InvalidRequestError
	}

	var user = &entities.Investor{}
	query := i.db.GormDB().
		WithContext(ctx).
		Model(user)
	query = filterInvestor(query, req)

	err := query.First(user).Error
	if err != nil {
		i.logger.Error("find user failed", info_log.ErrorToLogFields("details", err))
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, database.NotFoundError
		}
		return nil, err
	}

	return user, nil
}

func filterInvestor(db *gorm.DB, req query_params.GetUserParams) *gorm.DB {
	return db.Where("address = ?", req.Address)
}

func isValidParams(req query_params.GetUserParams) bool {
	return len(req.Address) > 0
}
