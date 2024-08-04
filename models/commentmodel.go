package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID          primitive.ObjectID `bson:"_id"`
	Post_ID     primitive.ObjectID `json:"post_id" validate:"required"`
	User_ID     primitive.ObjectID `json:"user_id" validate:"required"`
	Description *string            `json:"description" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
}
