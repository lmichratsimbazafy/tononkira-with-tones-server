package routes

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"lmich.com/tononkira/auth"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/http"
	"lmich.com/tononkira/mongodb"
)

type Router struct {
	Router *gin.Engine
}

func (r *Router) InitRoutes() {
	r.Router = gin.New()
	r.Router.Use(gin.Recovery()) // to recover gin automatically
	// r.Router.Use(middleware.JSONLogMiddleware()) // we'll define it later

	r.Router.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("authToken"))))

	// Create MongoDB service
	authService := auth.NewJWTService("mysecretkey", time.Hour*24)
	lyricsService := &mongodb.MongoLyricsService{Collection: config.GetCollections().LyricsModel}
	userService := &mongodb.MongoUserService{Collection: config.GetCollections().UserModel}
	authMiddleware := auth.NewAuthMiddleware(authService)
	// Handlers
	authHandler := &http.AuthHandler{UserService: userService, AuthService: authService}
	// userHandler := &http.Handler{UserService: userService}
	lyricsHandler := &http.LyricshHandler{LyricsService: lyricsService}

	r.Router.POST("/login", authHandler.Login)

	// r.Router.POST("/programs", program.Create)
	// r.Router.PUT("/programs/:id", program.Update)

	// r.Router.GET("/programs/:date", program.Details)
	// r.Router.GET("/authors/list", author.List)

	r.Router.GET("/lyrics/list", authMiddleware.AuthMiddleware(), lyricsHandler.List)
	r.Router.GET("/lyrics/list/authors/:id", authMiddleware.AuthMiddleware(), lyricsHandler.List)
}

func (r *Router) Run() {
	env := config.Getenv()
	r.Router.Run(":" + env.ApiPort)
}
