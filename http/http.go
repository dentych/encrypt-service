package http

import (
	"errors"
	"github.com/dentych/encrypt-service/api"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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
		filename := c.PostForm("filename")
		if filename == "" {
			c.String(http.StatusBadRequest, "filename is required")
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "failed to get form file")
			return
		}

		f, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to open file")
			return
		}
		defer f.Close()

		if err = api.SaveFile(filename, f); err != nil {
			if errors.Is(err, api.ErrFileExists) {
				c.String(http.StatusBadRequest, "file already exists.")
			} else {
				c.String(http.StatusInternalServerError, "failed to save file: %s", err)
			}
			return
		}

		c.Status(http.StatusOK)
	})
	h.g.GET("/file", func(c *gin.Context) {
		filename := c.Query("filename")
		if filename == "" {
			c.String(http.StatusBadRequest, "filename query parameter is required")
			return
		}

		reader, err := api.RetrieveFile(filename)
		if err != nil {
			c.String(http.StatusInternalServerError, "error retrieving file: %s", err)
			return
		}

		c.Status(http.StatusOK)
		_, err = io.Copy(c.Writer, reader)
		if err != nil {
			log.Printf("failed to copy to http writer: %s", err)
		}
	})

	return nil
}

func (h *Http) Start() error {
	err := h.g.Run(":8080")
	return err
}
