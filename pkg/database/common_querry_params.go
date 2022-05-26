package database

type CommonQueryParams struct {
	ID int64 `json:"id"`

	// by time
	CreatedAtFrom int64 `json:"created_at_from"`
	CreatedAtTo   int64 `json:"created_at_to"`

	// pagination
	Limit   int64 `json:"limit"`
	Offset  int64 `json:"offset"`
	SinceId int64 `json:"since_id"`
}
