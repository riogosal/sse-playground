package main

import (
	"github.com/gin-gonic/gin"
	"github.com/riogosal/sse-playground/sse"
)

type Hero struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PatchHeroRequest struct {
	Name string `json:"name"`
}

func main() {
	r := gin.New()

	hero := Hero{ID: 1, Name: "Batman"}
	signal := make(chan struct{})
	listeners := 0
	defer close(signal)

	r.GET("/hero", func(c *gin.Context) {
		c.JSON(200, hero)
	})

	r.GET("/hero/stream", func(c *gin.Context) {
		listeners++
		err_ch := make(chan error)

		defer func() {
			if listeners > 0 {
				listeners--
			}
		}()

		sse.OpenStream(c, signal, err_ch)
	})

	r.PATCH("/hero", func(c *gin.Context) {
		var req PatchHeroRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if req.Name == "" {
			c.JSON(400, gin.H{"error": "Hero name can't be empty"})
			return
		}

		if hero.Name == req.Name {
			c.JSON(200, hero)
			return
		}

		hero.Name = req.Name
		for i := 0; i < listeners; i++ {
			signal <- struct{}{}
		}

		c.JSON(200, hero)
	})

	r.Run(":8080")
}
