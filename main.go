package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Model for starwars
type Model struct {
	SetID    int    `json:"set-id"`
	Number   int    `json:"number"`
	Variant  int    `json:"variant"`
	Theme    string `json:"theme"`
	Subtheme string `json:"sub-theme"`
	Year     int    `json:"year"`
	Name     string `json:"name"`
	Minifigs int    `json:"minifigs"`
	Pieces   int    `json:"prices"`
	UKPrice  int    `json:"uk-price"`
	USPrice  int    `json:"us-price"`
	CAPrice  int    `json:"ca-price"`
	EUPrice  int    `json:"eu-price"`
	ImageURL string `json:"image-url"`
}

func main() {
	client := &http.Client{}
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		resp, err := client.Get("http://localhost:9200/")
		if err == nil {
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				c.Data(http.StatusOK, "", bytes)
			}
		}
		c.Status(http.StatusInternalServerError)
		c.Error(err)
	})

	router.POST("/write", func(c *gin.Context) {
		data, err := json.Marshal(Model{})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", "http://localhost:9200/hoges/hoge/test", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				c.Data(http.StatusOK, "application/json", bytes)
				return
			}
		}
		c.Status(http.StatusInternalServerError)
		c.Error(err)
		c.Abort()
	})

	router.Run(":8080")
}
