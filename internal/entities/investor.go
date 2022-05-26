package entities

type Investor struct {
	DefaultModel
	Address     string `json:"address" gorm:"size:512;uniqueIndex"`
	NetworkName string `json:"network_name" gorm:"size:32"`
	NetworkURL  string `json:"network_url" gorm:"size:512"`
	ChainID     int64  `json:"chain_id"`
	Symbol      string `json:"symbol" gorm:"size:16"`
	RefreshKey  string `json:"refresh_key" gorm:"size:512"`
}
