package src

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	service *Migrator

	eng *gin.Engine
	srv *http.Server
}

func NewServer(port int, cfg DBConfig) (*Server, error) {
	service, err := NewMigrator(cfg)
	if err != nil {
		return nil, err
	}

	eng := gin.Default()

	srv := &Server{
		eng: eng,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: eng,
		},
		service: service,
	}

	srv.registerHandlers()
	return srv, nil
}

func (srv *Server) Run() {
	if err := srv.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func (srv *Server) Close() {
	_ = srv.srv.Shutdown(context.Background())
}

func (srv *Server) registerHandlers() {
	router := srv.eng.Group("/v1")

	router.Use(AddCorrelationID)

	router.PUT("/create", srv.handleCreate())
	router.POST("/update", srv.handleUpdate())
	router.GET("/read", srv.handleRead())
	router.DELETE("/delete", srv.handleDelete())
}

func (srv *Server) handleCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m Model
		if err := c.ShouldBind(&m); err != nil {
			c.Status(http.StatusBadRequest)
		}
		created, err := srv.service.Create(c.Request.Context(), m)
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.JSON(200, created)
	}
}

func (srv *Server) handleUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m Model
		if err := c.ShouldBind(&m); err != nil {
			c.Status(http.StatusBadRequest)
		}
		if err := srv.service.Update(c.Request.Context(), m); err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.Status(200)
	}
}

func (srv *Server) handleRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		m, err := srv.service.Read(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.JSON(200, m)
	}
}

func (srv *Server) handleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m Model
		if err := c.ShouldBind(&m); err != nil {
			c.Status(http.StatusBadRequest)
		}
		created, err := srv.service.Create(c.Request.Context(), m)
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.JSON(200, created)
	}
}
