package dto

type CreateEquipmentRequest struct {
	Name              string `json:"name" validate:"required"`
	Category          string `json:"category"`
	QuantityTotal     int    `json:"quantity_total" validate:"min=0"`
	QuantityAvailable int    `json:"quantity_available" validate:"min=0"`
	Condition         string `json:"condition"`
	Notes             string `json:"notes"`
}

type UpdateEquipmentRequest struct {
	Name              *string `json:"name"`
	Category          *string `json:"category"`
	QuantityTotal     *int    `json:"quantity_total"`
	QuantityAvailable *int    `json:"quantity_available"`
	Condition         *string `json:"condition"`
	Notes             *string `json:"notes"`
}

type LoanEquipmentRequest struct {
	BorrowedBy uint   `json:"borrowed_by" validate:"required"`
	ActivityID *uint  `json:"activity_id"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
	DueDate    string `json:"due_date" validate:"required"`
}
