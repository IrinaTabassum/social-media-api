package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   error  `json:"error"`
}

type UserRegisterInput struct {
	First_name *string `json:"first_name" validate:"required,min=2,max=100"`
	Last_name  *string `json:"last_name" validate:"required,min=2,max=100"`
	Password   *string `json:"Password" validate:"required"`
	Email      *string `json:"email" validate:"email,required"`
	Phone      *string `json:"phone" validate:"required"`
}
type UserLoginInput struct {
	Password *string `json:"Password" validate:"required"`
	Email    *string `json:"email" validate:"email,required"`
}
type PostInput struct {
	Name        *string `json:"name" validate:"required"`
	Description *string `json:"description" validate:"required"`
}

type PostOutput struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `json:"name" validate:"required"`
	Description    string             `json:"description" validate:"required"`
	Comments       []Comment
	Total_Comments int                `json:"total_comments"`
	Likes          []Like
	Total_Likes    int                `json:"total_likes"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
	User_id        primitive.ObjectID `json:"user_id"`
}

type Poststac struct {
	ID             primitive.ObjectID `bson:"_id"`
	Name           string             `json:"name" validate:"required"`
	Description    string             `json:"description" validate:"required"`
	Comments_Count int                `json:"comment_count"`
	Likes_Count    int                `json:"like_count"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
	User_id        primitive.ObjectID `json:"user_id"`
}

type PostList struct {
	Posts []Poststac
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type CommentInput struct {
	Post_ID     primitive.ObjectID `json:"post_id" validate:"required"`
	Description *string            `json:"description" validate:"required"`
}

type CommentList struct {
	Comments []Comment
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type LikeInput struct {
	Post_ID primitive.ObjectID `json:"post_id" validate:"required"`
}

type LikeList struct {
	Likes []Like
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type CreateOutput struct {
	InsertedID string
}
