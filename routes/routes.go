package routes

import (
	"github.com/Atheer-Ganayem/Go-social-media-api/middlewares"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	authRoutes := server.Group("/", middlewares.IsAuth)

	server.POST("/login", login)
	server.POST("/signup", signup)

	// feed + search (users/posts)

	authRoutes.GET("/feed", getFeed)
	authRoutes.GET("/search", search)

	authRoutes.GET("/profile/:id", getProfile)
	authRoutes.PUT("/profile", updateProfile)
	authRoutes.POST("/follow/:id", follow) // follow/unfollow

	authRoutes.POST("/post", createPost)
	authRoutes.DELETE("/post/:id", deletePost)
	authRoutes.POST("/post/like/:id", likePost) // like/unlike

	authRoutes.POST("/comment/:postId", createComment)
	authRoutes.DELETE("/comment/:id", deleteComment)
}