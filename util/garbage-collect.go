package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/janlauber/autokueng-api/database"
	"github.com/janlauber/autokueng-api/models"
)

var (
	DataURL string = "https://data.autokueng.ch"
)

func GarbageCollect() {
	// TODO
	db := database.DBConn
	for {
		time.Sleep(time.Hour * 1)
		images := make(map[string][]string)
		var linkImages []string
		var memberImages []string
		var newsImages []string
		var serviceImages []string
		// get all images from models.Member.Image URLs
		db.Model(&models.Member{}).Pluck("image", &memberImages)
		db.Model(&models.Link{}).Pluck("image", &linkImages)
		db.Model(&models.News{}).Pluck("image", &newsImages)
		db.Model(&models.Service{}).Pluck("image", &serviceImages)

		// append all images to images map with key "activeImages"
		images["activeImages"] = append(images["activeImages"], memberImages...)
		images["activeImages"] = append(images["activeImages"], linkImages...)
		images["activeImages"] = append(images["activeImages"], newsImages...)
		images["activeImages"] = append(images["activeImages"], serviceImages...)

		// split images images string after "/images/"
		for _, image := range images["activeImages"] {
			if strings.Contains(image, "/images/") {
				images["activeImages"] = append(images["activeImages"], strings.Split(image, "/images/")[1])
			}
		}

		// POST request at localhost:9000/garbage-collect with images map
		jsonBytes, err := json.Marshal(images)
		if err != nil {
			log.Println(err)
		}
		req, err := http.NewRequest("POST", DataURL+"/garbage-collect", bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		client.Do(req)
	}

}
