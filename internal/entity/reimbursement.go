package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// Reimbursement represents an employee's reimbursement record
// swagger:model Reimbursement
type Reimbursement struct {
	// Unique identifier for the reimbursement record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	ID gorm.ULID `json:"id" gorm:"column:id;type:ulid;primaryKey"`

	// Reimbursement amount in currency units
	// example: 150000
	Amount int `json:"amount" gorm:"column:amount;type:integer;not null"`

	// Description of the reimbursement expense
	// example: "Transportation expenses for client meeting"
	Description string `json:"description" gorm:"column:description;type:text;not null"`

	// Timestamp when the reimbursement record was created
	// example: "2024-01-15T08:00:00Z"
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`

	// ID of the employee who created the reimbursement record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	CreatedBy gorm.ULID `json:"created_by" gorm:"column:created_by;type:ulid;not null"`

	// Timestamp when the reimbursement record was last updated
	// example: "2024-01-15T08:00:00Z"
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`

	// ID of the employee who last updated the reimbursement record
	// example: "01HXYZ123456789ABCDEFGHIJK"
	UpdatedBy *gorm.ULID `json:"updated_by" gorm:"column:updated_by;type:ulid"`

	// Relations
	// Employee who created the reimbursement record
	Creator *Employee `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	// Employee who last updated the reimbursement record
	Updater *Employee `json:"updater,omitempty" gorm:"foreignKey:UpdatedBy"`
}

// CreateReimbursementProps represents the properties needed to create a new reimbursement record
// swagger:model CreateReimbursementProps
type CreateReimbursementProps struct {
	// Reimbursement amount in currency units
	Amount int
	// Description of the reimbursement expense
	Description string
	// ID of the employee creating the reimbursement record
	CreatedBy gorm.ULID
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
