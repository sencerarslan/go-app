package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    *string            `json:"user_id"`
	Name      *string            `json:"name" validate:"required"`
	Logo      *string            `json:"logo" validate:"required"`
	Banner    *string            `json:"banner" validate:"required"`
	MenuGroup []MenuGroup        `json:"menu_groups"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
type MenuGroup struct {
	ID        primitive.ObjectID `bson:"_id"`
	MenuID    *string            `json:"menu_id" validate:"required"`
	Name      *string            `json:"name" validate:"required"`
	MenuItem  []MenuItem         `json:"menu_items"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	GroupID     *string            `json:"group_id"`
	Name        *string            `json:"name" validate:"required"`
	Price       float64            `json:"price" validate:"required"`
	Description *string            `json:"description" validate:"required"`
	ImageURL    *string            `json:"image_url" validate:"required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
