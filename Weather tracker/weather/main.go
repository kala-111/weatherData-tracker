package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type weatherData struct {
	Name string `json:"name"`
	gorm.Model
	Main struct {
		kelvin float64 `json:"temp"`
	} `json:"main"`
}

var wd weatherData

func GetUsers(db *gorm.DB, User *[]weatherData) (err error) {
	err = db.Find(wd).Error
	if err != nil {
		return err
	}
	return nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from GoLang!\n"))
}

func query(city string) (weatherData, error) {

	resp, err := http.Get("http://api.Open.meteo.com/v1/forecast?" + "&=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return wd, nil

}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", nil)
}
