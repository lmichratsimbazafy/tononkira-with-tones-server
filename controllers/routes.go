package controllers

import (
	"github.com/gin-gonic/gin"
	"lmich.com/tononkira/config"
	"lmich.com/tononkira/controllers/author"
	"lmich.com/tononkira/controllers/lyrics"
)

type Router struct {
	router *gin.Engine
}

func (r *Router) Init() {
	r.router = gin.Default()
	r.router.GET("/authors/list", author.List)
	r.router.GET("/lyrics/list", lyrics.List)
	r.router.GET("/lyrics/list/authors/:id", lyrics.List)
}

func (r *Router) Run() {
	env := config.Getenv()
	r.router.Run(":" + env.ApiPort)
}
