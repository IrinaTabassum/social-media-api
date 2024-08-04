package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
    ID            primitive.ObjectID `bson:"_id"`
    Name          *string            `json:"name" validate:"required"`
    Description   *string            `json:"description" validate:"required"`
    Created_at    time.Time          `json:"created_at"`
    Updated_at    time.Time          `json:"updated_at"`
    User_id       primitive.ObjectID `json:"user_id"`
}	

