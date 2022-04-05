package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Http struct {
	g *gin.Engine
}

func CreateHttp() *Http {
	return &Http{}
}

func (h *Http) Setup() error {
	h.g = gin.Default()
	h.g.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "failed to get form file")
			return
		}

		f, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to open file")
		}
		defer f.Close()
		c.Status(http.StatusOK)
	})

	return nil
}

func (h *Http) Start() error {
	err := h.g.Run(":8080")
	return err
}
