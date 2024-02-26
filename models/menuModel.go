package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID         primitive.ObjectID `bson:"_id"`
	Menu_id    *string            `json:"menu_id"`
	User_id    *string            `json:"user_id"`
	Name       *string            `json:"name" validate:"required"`
	MenuItem   []MenuItem         `json:"menu_items"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
type MenuItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Menu_id     *string            `json:"menu_id"`
	Name        *string            `json:"name" validate:"required"`
	Price       float64            `json:"price" validate:"required"`
	Description *string            `json:"description" validate:"required"`
	ImageURL    *string            `json:"image_url" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
