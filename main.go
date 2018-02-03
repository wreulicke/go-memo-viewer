package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get("http://localhost:9200/")
		if err == nil {
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				w.Write(bytes)
				return
			}
			return
		}
		http.Error(w, resp.Status, resp.StatusCode)
	})
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(Model{})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:9200/hoges/hoge/test", bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			bytes, err := ioutil.ReadAll(resp.Body)
			fmt.Println(resp.Status, string(bytes))
			if err == nil {
				w.Write(bytes)
				return
			}
			return
		}
		http.Error(w, resp.Status, resp.StatusCode)
	})
	http.ListenAndServe(":8080", nil)
}
