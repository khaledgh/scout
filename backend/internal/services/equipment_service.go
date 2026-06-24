package services

import (
	"context"
	"kashfi/internal/models"
	"time"

	"gorm.io/gorm"
)

type EquipmentService struct {
	db *gorm.DB
}

func NewEquipmentService(db *gorm.DB) *EquipmentService {
	return &EquipmentService{db: db}
}

func (s *EquipmentService) List(ctx context.Context) ([]models.Equipment, error) {
	var items []models.Equipment
	err := s.db.WithContext(ctx).Order("name").Find(&items).Error
	return items, err
}

func (s *EquipmentService) Get(ctx context.Context, id uint) (*models.Equipment, error) {
	var item models.Equipment
	if err := s.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, ErrNotFound
	}
	return &item, nil
}

type CreateEquipmentInput struct {
	Name              string
	Category          string
	QuantityTotal     int
	QuantityAvailable int
	Condition         string
	Notes             string
}

func (s *EquipmentService) Create(ctx context.Context, in CreateEquipmentInput) (*models.Equipment, error) {
	item := &models.Equipment{
		Name:              in.Name,
		Category:          in.Category,
		QuantityTotal:     in.QuantityTotal,
		QuantityAvailable: in.QuantityAvailable,
		Condition:         in.Condition,
		Notes:             in.Notes,
	}
	if err := s.db.WithContext(ctx).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

type UpdateEquipmentInput struct {
	Name              *string
	Category          *string
	QuantityTotal     *int
	QuantityAvailable *int
	Condition         *string
	Notes             *string
}

func (s *EquipmentService) Update(ctx context.Context, id uint, in UpdateEquipmentInput) (*models.Equipment, error) {
	var item models.Equipment
	if err := s.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, ErrNotFound
	}
	updates := map[string]interface{}{}
	if in.Name != nil { updates["name"] = *in.Name }
	if in.Category != nil { updates["category"] = *in.Category }
	if in.QuantityTotal != nil { updates["quantity_total"] = *in.QuantityTotal }
	if in.QuantityAvailable != nil { updates["quantity_available"] = *in.QuantityAvailable }
	if in.Condition != nil { updates["condition"] = *in.Condition }
	if in.Notes != nil { updates["notes"] = *in.Notes }

	if err := s.db.WithContext(ctx).Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *EquipmentService) Delete(ctx context.Context, id uint) error {
	result := s.db.WithContext(ctx).Delete(&models.Equipment{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (s *EquipmentService) Loans(ctx context.Context, equipmentID uint) ([]models.EquipmentLoan, error) {
	var loans []models.EquipmentLoan
	err := s.db.WithContext(ctx).Where("equipment_id = ?", equipmentID).
		Preload("Borrower").Preload("Activity").Order("created_at DESC").Find(&loans).Error
	return loans, err
}

type LoanInput struct {
	EquipmentID uint
	BorrowedBy  uint
	ActivityID  *uint
	Quantity    int
	DueDate     time.Time
}

// Loan creates a loan and decrements available quantity inside a transaction.
func (s *EquipmentService) Loan(ctx context.Context, in LoanInput) (*models.EquipmentLoan, error) {
	var loan *models.EquipmentLoan
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item models.Equipment
		if err := tx.First(&item, in.EquipmentID).Error; err != nil {
			return ErrNotFound
		}
		if item.QuantityAvailable < in.Quantity {
			return ErrBadRequest
		}
		l := &models.EquipmentLoan{
			EquipmentID: in.EquipmentID,
			BorrowedBy:  in.BorrowedBy,
			ActivityID:  in.ActivityID,
			Quantity:    in.Quantity,
			DueDate:     in.DueDate,
		}
		if err := tx.Create(l).Error; err != nil {
			return err
		}
		if err := tx.Model(&item).Update("quantity_available", item.QuantityAvailable-in.Quantity).Error; err != nil {
			return err
		}
		loan = l
		return nil
	})
	return loan, err
}

// ReturnLoan marks a loan returned and restores available quantity.
func (s *EquipmentService) ReturnLoan(ctx context.Context, loanID uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var loan models.EquipmentLoan
		if err := tx.First(&loan, loanID).Error; err != nil {
			return ErrNotFound
		}
		if loan.ReturnedAt != nil {
			return ErrBadRequest
		}
		now := time.Now()
		if err := tx.Model(&loan).Update("returned_at", &now).Error; err != nil {
			return err
		}
		var item models.Equipment
		if err := tx.First(&item, loan.EquipmentID).Error; err != nil {
			return err
		}
		return tx.Model(&item).Update("quantity_available", item.QuantityAvailable+loan.Quantity).Error
	})
}
