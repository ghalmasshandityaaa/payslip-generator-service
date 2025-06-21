package repository

import (
	"payslip-generator-service/internal/model"
	ulid "payslip-generator-service/pkg/database/gorm"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Debug().Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Debug().Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Debug().Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Debug().Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id ulid.ULID) error {
	return db.Debug().Where("id = ?", id).Take(entity).Error
}

func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T) error {
	return db.Debug().Find(entities).Error
}

func (r *Repository[T]) FindAllWithPagination(db *gorm.DB, pagination *model.PaginationOptions) ([]T, int64, error) {
	var entities []T
	query := db
	if pagination.Filter != nil {
		query = query.Scopes(*pagination.Filter)
	}
	if len(pagination.Order) > 0 {
		for _, order := range pagination.Order {
			query = query.Order(order.Column + " " + string(order.Direction))
		}
	}
	if err := query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := db
	if pagination.Filter != nil {
		countQuery = countQuery.Scopes(*pagination.Filter)
	}
	if err := countQuery.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}
