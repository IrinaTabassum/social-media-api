package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID         primitive.ObjectID `bson:"_id"`
	User_ID    primitive.ObjectID `json:"user_id" validate:"required"`
	Post_ID    primitive.ObjectID `json:"post_id" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
}
