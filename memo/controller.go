package memo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// Model for starwars
type Memo struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type MemoSearchRequest struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
}

type Controller struct {
	http http.Client
	db   *sql.DB
}

func (controller *Controller) post(c *gin.Context) {
	model := &Memo{}
	err := c.BindJSON(model)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id := uuid.NewV4()
	_, err = controller.db.Exec("insert into memo (id, title, text) values (?, ?, ?)", id.String(), model.Text, model.Title)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(model)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:9200/memos/memo", bytes.NewReader(data))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	resp, err := controller.http.Do(req)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	c.Data(http.StatusOK, "application/json", bytes)
}

func (controller *Controller) query(c *gin.Context) {
	request := &MemoSearchRequest{}
	err := c.BindJSON(request)

	if err != nil {
		c.String(http.StatusBadRequest, `Required JSON Request`)
		return
	}

	if request.Title == nil && request.Text == nil {
		c.String(http.StatusBadRequest, `Required "text" or "title"`)
	}

	query := []string{}
	if request.Text != nil {
		query = append(query, fmt.Sprintf(`{ "match": { "text": %q } }`, *request.Text))
	}
	if request.Title != nil {
		query = append(query, fmt.Sprintf(`{ "match": { "title": %q } }`, *request.Title))
	}

	queryStr := strings.Join(query, ",")

	rs := fmt.Sprintf(`{
		"query": {
			"bool": {
				"must": [
					%s
				]
			}
		}
	}`, queryStr)
	req, err := http.NewRequest("GET", "http://localhost:9200/memos/_search", bytes.NewReader([]byte(rs)))
	req.Header.Add("Content-Type", "application/json")

	resp, err := controller.http.Do(req)
	if err != nil {
		c.Error(err)
		c.String(http.StatusServiceUnavailable, "Error")
		return
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Error")
		return
	}

	c.Data(http.StatusOK, "application/json", bytes)
}

// Route registers routes.
func Route(group *gin.RouterGroup, conn *sql.DB) {
	controller := Controller{
		http: http.Client{},
		db:   conn,
	}
	group.POST("", controller.post)
	group.GET("/_search", controller.query)
}
