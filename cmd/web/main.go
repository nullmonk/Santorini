package main

import (
	"fmt"
	"net/http"
	"santorini/santorini"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/scripts", func(c *gin.Context) {
		c.JSON(http.StatusOK, santorini.Storage.Scripts())
	})
	r.POST("/script", func(c *gin.Context) {
		// Parse JSON
		var json struct {
			Password string `json:"password" binding:""`
			Name     string `json:"name" binding:"required"`
			Contents string `json:"contents" binding:"required"`
		}
		err := c.Bind(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}
		err = santorini.Storage.SaveScript(json.Password, &santorini.Script{
			Name:     json.Name,
			Contents: json.Contents,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": nil})
	})

	// Play a single round
	r.POST("/play", func(c *gin.Context) {
		var json struct {
			Players []string `json:"players" binding:"required"`
		}
		err := c.Bind(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}

		players := make([]*santorini.Player, 0, len(json.Players))
		for _, p := range json.Players {
			s, err := santorini.Storage.LoadScript(p)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("error loading player '%s': %s", p, err)})
				return
			}
			players = append(players, &santorini.Player{
				Source: s.Contents,
				Name:   s.Name,
			})
		}
		g := santorini.NewGame(1, nil, players...)
		err = g.Finish()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("error playing game: %s", err)})
			return
		}
		// TODO Save game logs here
		c.JSON(http.StatusOK, gin.H{"status": "ok", "log": g.GetTextLog()})
	})

	// Simulate (POST), trigger a simulation and open a websocket for game results

	r.Run(":8080")
}
