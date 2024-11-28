package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var publicUrl = "frontend/public/"

var apiAddresses = map[string]string{
	"artists": "https://groupietrackers.herokuapp.com/api/artists",
}

type artistsData struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ResultPageData struct {
	Result          string
	Color           string
	Align           string
	BackgroundColor string
}

type ErrorPageData struct {
	Name       string
	Code       string
	CodeNumber int
	Info       string
}

var PredefinedErrors = map[string]ErrorPageData{
	"BadRequestError": {
		Name:       "BadRequestError",
		Code:       strconv.Itoa(http.StatusBadRequest),
		CodeNumber: http.StatusBadRequest,
		Info:       "Bad request",
	},
	"NotFoundError": {
		Name:       "NotFoundError",
		Code:       strconv.Itoa(http.StatusNotFound),
		CodeNumber: http.StatusNotFound,
		Info:       "Page not found",
	},
	"MethodNotAllowedError": {
		Name:       "MethodNotAllowedError",
		Code:       strconv.Itoa(http.StatusMethodNotAllowed),
		CodeNumber: http.StatusMethodNotAllowed,
		Info:       "Method not allowed",
	},
	"InternalServerError": {
		Name:       "InternalServerError",
		Code:       strconv.Itoa(http.StatusInternalServerError),
		CodeNumber: http.StatusInternalServerError,
		Info:       "Internal server error",
	},
}

var (
	BadRequestError       = PredefinedErrors["BadRequestError"]
	NotFoundError         = PredefinedErrors["NotFoundError"]
	MethodNotAllowedError = PredefinedErrors["MethodNotAllowedError"]
	InternalServerError   = PredefinedErrors["InternalServerError"]
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			// If the URL is not exactly "/", respond with 404
			handleErrorPage(w, r, NotFoundError)
			return
		}

		req, err := http.NewRequest("GET", apiAddresses["artists"], nil)
		if err != nil {
			fmt.Print(err.Error())
		}
		// req.Header.Add("x-rapidapi-key", "YOU_API_KEY")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Print(err.Error())
		}
		defer res.Body.Close()
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			fmt.Print(err.Error())
		}
		// fmt.Println(string(body))

		var data_obj []artistsData
		jsonErr := json.Unmarshal(body, &data_obj)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		fmt.Println("\n________________________________________")

		for _, data := range data_obj {
			fmt.Println("Id: ", data.Id)
			fmt.Println("Image: ", data.Image)
			fmt.Println("Name: ", data.Name)
			fmt.Println("CreationDate: ", data.CreationDate)
			fmt.Println("FirstAlbum: ", data.FirstAlbum)
			fmt.Println("Locations: ", data.Locations)
			fmt.Println("ConcertDates: ", data.ConcertDates)
			fmt.Println("Relations: ", data.Relations)
			for _, value := range data.Members {
				fmt.Println("Members: ")
				fmt.Print(value)
			}
			fmt.Println("\n________________________________________\n")
		}

		tmpl, err := template.ParseFiles(publicUrl + "index.html")

		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}
		tmpl.Execute(w, data_obj)
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleErrorPage(w http.ResponseWriter, r *http.Request, errorType ErrorPageData) {
	tmpl, err := template.ParseFiles("frontend/errors/error.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(errorType.CodeNumber)
	tmpl.Execute(w, errorType)
}

func main() {
	http.Handle("/static/", http.FileServer(http.Dir("./frontend/public/")))
	http.HandleFunc("/", handleIndex)
	// http.HandleFunc("/ascii-web", handleAsciiWeb)
	// Start the server on port 8082
	fmt.Println("Starting server on 0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
