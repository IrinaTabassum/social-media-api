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
)
var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")


//CreateUser is the api used to tget a single user
func CreatePost() gin.HandlerFunc {
    return func(c *gin.Context) {
		uid, exists := c.Get("uid")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		UID, err := primitive.ObjectIDFromHex(uid.(string))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UID"})
            return
        }
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var post models.Post
        if err := c.BindJSON(&post); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        validationErr := validate.Struct(post)
        if validationErr != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
		
		post.User_id = UID;
        post.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        post.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        post.ID = primitive.NewObjectID()
        

        resultInsertionNumber, insertErr := postCollection.InsertOne(ctx, post)
        if insertErr != nil {
            msg := fmt.Sprintf("Post item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }
        defer cancel()

        c.JSON(http.StatusOK, resultInsertionNumber)

    }
}

// The function `GetPostByID` retrieves a post by its ID along with associated comments and likes using
// MongoDB aggregation pipeline in a Go application.
func GetPostByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        postID := c.Param("id")
        objID, err := primitive.ObjectIDFromHex(postID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
            return
        }

        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        // Aggregation pipeline
        pipeline := mongo.Pipeline{
            {{Key: "$match", Value: bson.D{{"_id", objID}}}},
            {{Key: "$lookup", Value: bson.D{{"from", "comment"}, {"localField", "_id"}, {"foreignField", "post_id"},{"as", "comments"},}}},
            {{Key: "$lookup", Value: bson.D{{"from", "like"},{"localField", "_id"},{"foreignField", "post_id"},{"as", "likes"},}}},
            {{Key: "$addFields", Value: bson.D{
                {"total_comments", bson.D{{"$size", "$comments"}}},
                {"total_likes", bson.D{{"$size", "$likes"}}},
            }}},
        }

        cursor, err := postCollection.Aggregate(ctx, pipeline)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching post"})
            return
        }

        var posts []bson.M
        if err = cursor.All(ctx, &posts); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding post data"})
            return
        }
		fmt.Println(posts[0])

        c.JSON(http.StatusOK, posts[0])
    }
}

// The function `UpdatePost` in Go handles updating a post in a web application, including
// authorization checks and database operations.
func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
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

		var post models.Post
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find the post and check if the user is authorized to update it
		var existingPost models.Post
		filter := bson.M{"_id": objID, "user_id": UID}
		if err := postCollection.FindOne(context.Background(), filter).Decode(&existingPost); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		post.Updated_at = time.Now()
		update := bson.M{
			"$set": bson.M{
				"name":       post.Name,
				"description": post.Description,
				"updated_at":  post.Updated_at,
			},
		}

		_, err = postCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
	}
}

// DeletePost deletes a post by ID
func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
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


		filter := bson.M{"_id": objID, "user_id": UID}
		result, err := postCollection.DeleteOne(context.Background(), filter)
		if err != nil || result.DeletedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post or not authorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}

// The ListPosts function retrieves paginated post data with comments and likes aggregation using
// MongoDB aggregation pipeline in a Go application.
func ListPosts() gin.HandlerFunc {
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

        // Calculate skip for pagination
        skip := int64((page - 1) * limit)


        // Aggregation pipeline
        pipeline := mongo.Pipeline{
            {{Key: "$lookup", Value: bson.D{{"from", "comment"}, {"localField", "_id"}, {"foreignField", "post_id"}, {"as", "comments"}}}},
            {{Key: "$lookup", Value: bson.D{{"from", "like"}, {"localField", "_id"}, {"foreignField", "post_id"}, {"as", "likes"}}}},
            {{Key: "$project", Value: bson.D{
                {"name", 1},
                {"description", 1},
                {"created_at", 1},
                {"updated_at", 1},
                {"user_id", 1},
                {"like_count", bson.D{{"$size", "$likes"}}},
                {"comment_count", bson.D{{"$size", "$comments"}}},
            }}},
            {{"$skip", skip}},
            {{"$limit", int64(limit)}},
        }

        cursor, err := postCollection.Aggregate(ctx, pipeline)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
            return
        }

        var posts []bson.M
        if err = cursor.All(ctx, &posts); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding posts"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "data": posts,
            "page": page,
            "limit": limit,
        })
    }
}





