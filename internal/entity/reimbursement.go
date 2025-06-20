package entity

import (
	"time"

	"payslip-generator-service/pkg/database/gorm"

	"github.com/oklog/ulid/v2"
)

// ReimbursementStatus represents the status of a reimbursement
type ReimbursementStatus string

const (
	ReimbursementStatusPending  ReimbursementStatus = "pending"
	ReimbursementStatusApproved ReimbursementStatus = "approved"
	ReimbursementStatusRejected ReimbursementStatus = "rejected"
)

// Reimbursement model
type Reimbursement struct {
	ID          gorm.ULID           `json:"id" gorm:"column:id;type:ulid;primaryKey"`
	Amount      int                 `json:"amount" gorm:"column:amount;type:integer;not null"`
	Description string              `json:"description" gorm:"column:description;type:text;not null"`
	Status      ReimbursementStatus `json:"status" gorm:"column:status;type:reimbursement_status;not null;default:'pending'"`
	CreatedAt   time.Time           `json:"created_at" gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	CreatedBy   gorm.ULID           `json:"created_by" gorm:"column:created_by;type:ulid;not null"`
	UpdatedAt   *time.Time          `json:"updated_at" gorm:"column:updated_at;type:timestamp with time zone"`
	UpdatedBy   *gorm.ULID          `json:"updated_by" gorm:"column:updated_by;type:ulid"`

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
		Status:      ReimbursementStatusPending,
		CreatedAt:   time.Now(),
		CreatedBy:   gorm.ULID(props.CreatedBy),
	}
}

func (r *Reimbursement) TableName() string {
	return "reimbursement"
}

// IsPending checks if the reimbursement is in pending status
func (r *Reimbursement) IsPending() bool {
	return r.Status == ReimbursementStatusPending
}

// IsApproved checks if the reimbursement is approved
func (r *Reimbursement) IsApproved() bool {
	return r.Status == ReimbursementStatusApproved
}

// IsRejected checks if the reimbursement is rejected
func (r *Reimbursement) IsRejected() bool {
	return r.Status == ReimbursementStatusRejected
}

// Approve approves the reimbursement
func (r *Reimbursement) Approve(approvedBy gorm.ULID) {
	now := time.Now()
	r.Status = ReimbursementStatusApproved
	r.UpdatedAt = &now
	updatedBy := gorm.ULID(approvedBy)
	r.UpdatedBy = &updatedBy
}

// Reject rejects the reimbursement
func (r *Reimbursement) Reject(rejectedBy gorm.ULID) {
	now := time.Now()
	r.Status = ReimbursementStatusRejected
	r.UpdatedAt = &now
	updatedBy := gorm.ULID(rejectedBy)
	r.UpdatedBy = &updatedBy
}
