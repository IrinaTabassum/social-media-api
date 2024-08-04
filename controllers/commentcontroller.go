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
var commentCollection *mongo.Collection = database.OpenCollection(database.Client, "comment")


//CreateUser is the api used to tget a single user
func CreateComment() gin.HandlerFunc {
    return func(c *gin.Context) {
		uid, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var comment models.Comment
        if err := c.BindJSON(&comment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
		UID, err := primitive.ObjectIDFromHex(uid.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID"})
			return
		}
		comment.User_ID = UID
		
        comment.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        comment.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        comment.ID = primitive.NewObjectID()
		fmt.Println(comment)
		fmt.Println(comment.Post_ID)

        validationErr := validate.Struct(comment)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
		
        resultInsertionNumber, insertErr := commentCollection.InsertOne(ctx, comment)
        if insertErr != nil {
            msg := fmt.Sprintf("comment item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }
        defer cancel()

        c.JSON(http.StatusOK, resultInsertionNumber)

    }
}

// The function `GetCommentByID` retrieves a comment by its ID from a MongoDB collection and returns it
// as a JSON response using Gin framework in Go.
func GetCommentByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        commentID := c.Param("id")
        objectID, err := primitive.ObjectIDFromHex(commentID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
            return
        }
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        var comment models.Comment
        err = commentCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&comment)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Comment not found"})
            return
        }

        c.JSON(http.StatusOK, comment)
    }
}


// The UpdateComment function in Go handles updating a comment in a Gin framework application by
// validating the comment ID, user authorization, binding JSON data, updating the comment

func UpdateComment() gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID"})
			return
		}

        var updatedComment models.Comment
        if err := c.BindJSON(&updatedComment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        updatedComment.Updated_at = time.Now()

        filter := bson.M{"_id": objectID, "user_id": UID}
		update := bson.M{
			"$set": bson.M{
				"description": updatedComment.Description,
				"updated_at":  updatedComment.Updated_at,
			},
		}
        
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        result, err := commentCollection.UpdateOne(ctx, filter, update)
        if err != nil || result.MatchedCount ==0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Comment could not be updated"})
            return
        }
        

        c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
    }
}

// DeleteComment deletes a post by ID
func DeleteComment() gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid UID"})
			return
		}

		filter := bson.M{"_id": objectID, "user_id": UID}
		result, err := commentCollection.DeleteOne(context.Background(), filter)
		if err != nil || result.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment or not authorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	}
}

func GetCommentList() gin.HandlerFunc {
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

        var comments []models.Comment

        cursor, err := commentCollection.Find(ctx, bson.M{}, findOptions)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching comments"})
            return
        }

        if err = cursor.All(ctx, &comments); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding comments"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "data": comments,
            "page": page,
            "limit": limit,
        })
    }
}





