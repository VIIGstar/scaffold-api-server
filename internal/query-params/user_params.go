package query_params

import "scaffold-api-server/pkg/database"

type GetUserParams struct {
	database.CommonQueryParams
	Address string `json:"address"`
}
