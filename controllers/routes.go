package controllers

import (
	"github.com/gin-gonic/gin"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/controllers/author"
	"lmich.com/tononkira/controllers/lyrics"
	"lmich.com/tononkira/controllers/program"
)

type Router struct {
	router *gin.Engine
}

func (r *Router) Init() {
	r.router = gin.Default()
	r.router.POST("/programs", program.Create)
	r.router.PUT("/programs/:id", program.Update)
	r.router.GET("/programs/:date", program.Details)
	r.router.GET("/authors/list", author.List)
	r.router.GET("/lyrics/list", lyrics.List)
	r.router.GET("/lyrics/list/authors/:id", lyrics.List)
}

func (r *Router) Run() {
	env := config.Getenv()
	r.router.Run(":" + env.ApiPort)
}
