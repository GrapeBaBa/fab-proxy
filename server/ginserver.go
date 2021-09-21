package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGin(addr string) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	rg := router.Group("/api")
	rg.GET("/test", func(c *gin.Context) {
		//time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return srv
}
