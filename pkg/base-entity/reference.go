package base_entity

import "time"

// Reference information to refer other entity from 3rd parties/agencies
type Reference struct {
	RefId        int64     `json:"ref_id" gorm:"uniqueIndex"`
	RefType      string    `json:"ref_type" gorm:"type:varchar"`
	RefCreatedAt time.Time `json:"ref_created_at"`
	LinkedAt     time.Time `json:"linked_at"`
}
