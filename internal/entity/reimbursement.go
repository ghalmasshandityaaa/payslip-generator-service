package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Reimbursement model
type Reimbursement struct {
	ID          gorm.ULID  `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	Amount      int        `json:"amount" gorm:"column:amount;type:integer;not null"`
	Description string     `json:"description" gorm:"column:description;type:text;not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy   gorm.ULID  `json:"created_by" gorm:"column:created_by;type:ulid;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
	UpdatedBy   *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

type CreateReimbursementProps struct {
	Amount      int
	Description string
	CreatedBy   gorm.ULID
}

func NewReimbursement(props *CreateReimbursementProps) *Reimbursement {
	return &Reimbursement{
		ID:          gorm.ULID(ulid.Make()),
		Amount:      props.Amount,
		Description: props.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   gorm.ULID(props.CreatedBy),
	}
}

func (r *Reimbursement) TableName() string {
	return "reimbursement"
}
