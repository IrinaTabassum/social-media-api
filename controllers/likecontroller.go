package controller

import (
	"context"
	"fmt"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"social-media-api/database"
	"social-media-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var likeCollection *mongo.Collection = database.OpenCollection(database.Client, "like")

// CreateLike is the API used to create a new like
// @Summary Create a new Like
// @Description This endpoint allows a user to create a new like for a post.
// @Tags Like
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param like body models.Like true "Like Data"
// @Success 200 {object} models.CreateOutput
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /likes [post]
func CreateLike() gin.HandlerFunc {
    return func(c *gin.Context) {
		uid, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var like models.Like
        if err := c.BindJSON(&like); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
		UID, err := primitive.ObjectIDFromHex(uid.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID"})
			return
		}
		like.User_ID = UID;
        like.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        like.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        like.ID = primitive.NewObjectID()

        validationErr := validate.Struct(like)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
		
        resultInsertionNumber, insertErr := likeCollection.InsertOne(ctx, like)
        if insertErr != nil {
            msg := fmt.Sprintf("like item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }
        defer cancel()

        c.JSON(http.StatusOK, resultInsertionNumber)

    }
}
// GetLikeByID is the API used to fetch a single like by ID
// @Summary Get a Like by ID
// @Description This endpoint retrieves a like by its ID.
// @Tags Like
// @Accept json
// @Produce json
// @Param id path string true "Like ID"
// @Success 200 {object} models.Like
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /likes/{id} [get]
func GetLikeByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        likeID := c.Param("id")
        objectID, err := primitive.ObjectIDFromHex(likeID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
            return
        }
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        var like models.Like
        err = likeCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&like)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "like not found"})
            return
        }
        c.JSON(http.StatusOK, like)
    }
}
// DeleteLike deletes a like by ID
// @Summary Delete a Like
// @Description This endpoint allows a user to delete a like by ID.
// @Tags Like
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param id path string true "Like ID"
// @Success 200 {string} Like deleted successfully
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 401 {object} models.Error "Unauthorized"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /likes/{id} [delete]
func DeleteLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("id")
        objectID, err := primitive.ObjectIDFromHex(commentID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
            return
        }

		uid, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
        UID, err := primitive.ObjectIDFromHex(uid.(string))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post UID"})
            return
        }

		filter := bson.M{"_id": objectID, "user_id": UID}
		result, err := likeCollection.DeleteOne(context.Background(), filter)
		if err != nil || result.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment or not authorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
	}
}

// GetLikeList is the API used to get a list of likes with pagination
// @Summary Get a list of Likes
// @Description This endpoint retrieves a paginated list of likes.
// @Tags Like
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of likes per page" default(10)
// @Success 200 {object} models.LikeList
// @Failure 400 {object} models.Error "Invalid pagination parameters"
// @Failure 500 {object} models.Error "Internal server error"
// @Router /likes [GET]
func GetLikeList() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get pagination parameters from query string
        page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
        if err != nil || page < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
            return
        }

        limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
        if err != nil || limit < 1 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()

        // Define options for pagination
        findOptions := options.Find()
        findOptions.SetSkip(int64((page - 1) * limit))
        findOptions.SetLimit(int64(limit))

        var likes []models.Like

        cursor, err := likeCollection.Find(ctx, bson.M{}, findOptions)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching likes"})
            return
        }

        if err = cursor.All(ctx, &likes); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding likes"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "data": likes,
            "page": page,
            "limit": limit,
        })
    }
}







