package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
        "os"
	"time"
)


var influxURL string
var influxURI string

type StandardError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Appwatcher struct {
	App struct {
		Name string `json:"name"`
	} `json:"app"`
	Space struct {
		Name string `json:"name"`
	} `json:"space"`
	Dynos []struct {
		Dyno string `json:"dyno"`
		Type string `json:"type"`
	} `json:"dynos"`
	Key         string    `json:"key"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	Restarts    int       `json:"restarts"`
	CrashedAt   time.Time `json:"crashed_at"`
}

type Annotation struct {
	App   string
	Title string
	Text  string
	Tags  string
}

func setenv() {
    influxURL=os.Getenv("INFLUX_URL")
    influxURI=os.Getenv("INFLUX_URI")
}
func main() {
	setenv()
	api := gin.Default()
	api.POST("/v1/appwatcher", receive_appwatcher)
	api.Run(":"+os.Getenv("PORT"))

}

func receive_appwatcher(c *gin.Context) {
	var json Appwatcher
	err := c.BindJSON(&json)
	if err != nil {
		fmt.Println(err)
	}
	var annotation Annotation
     for _, element := range json.Dynos {
	annotation.App = json.Key
	annotation.Title = json.Action
	annotation.Text = json.Description
	annotation.Tags = json.Code + "," + json.Space.Name + "," + json.App.Name + "," + element.Type+"."+element.Dyno
	sendAnnotation(annotation)
      }
	c.JSON(200, nil)
}

func sendAnnotation(annotation Annotation) {
	client := http.Client{}
	data := "events  title=\"" + annotation.Title + "\",text=\"" + annotation.Text + "\",tags=\"" + annotation.Tags + "\",app=\"" + annotation.App + "\""
	databytes := []byte(data)

	req, err := http.NewRequest("POST",influxURL+influxURI, bytes.NewBuffer(databytes))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}
