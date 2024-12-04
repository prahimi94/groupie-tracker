package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var publicUrl = "frontend/public/"

var apiUrls = map[string]string{
	"base": "https://groupietrackers.herokuapp.com/api",
}

type ArtistsData struct {
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

type LocationsDataLevel2 struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationsDataLevel1 struct {
	Index []LocationsDataLevel2 `json:"index"`
}

type DatesDataLevel2 struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesDataLevel1 struct {
	Index []DatesDataLevel2 `json:"index"`
}

type RelationsDataLevel2 struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationsDataLevel1 struct {
	Index []RelationsDataLevel2 `json:"index"`
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

func sendGetRequest(url string, data_obj interface{}) {
	req, err := http.NewRequest("GET", url, nil)
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

	jsonErr := json.Unmarshal(body, &data_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func getApiUrls() {
	type ApiInfo struct {
		Artists   string `json:"artists"`
		Locations string `json:"locations"`
		Dates     string `json:"dates"`
		Relations string `json:"relation"`
	}

	var data_obj ApiInfo
	sendGetRequest(apiUrls["base"], &data_obj)

	apiUrls["artists"] = data_obj.Artists
	apiUrls["dates"] = data_obj.Dates
	apiUrls["locations"] = data_obj.Locations
	apiUrls["relations"] = data_obj.Relations
}

func generateUrl(path string, desiredUrl string) (string, string, string) {
	var url string
	if path == "/"+desiredUrl {
		url = apiUrls[desiredUrl]
		return url, "", ""
	} else if strings.HasPrefix(path, "/"+desiredUrl+"/") {
		id := strings.TrimPrefix(path, "/"+desiredUrl+"/")
		url = apiUrls[desiredUrl+"s"]
		return url, id, ""
	} else {
		return "", "", "not found"
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/" {
			// If the URL is not exactly "/", respond with 404
			handleErrorPage(w, r, NotFoundError)
			return
		}

		var data_obj []ArtistsData
		sendGetRequest(apiUrls["artists"], &data_obj)

		tmpl, err := template.ParseFiles(
			publicUrl+"index.html",
			publicUrl+"header.html",
			publicUrl+"artists.html",
		)

		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}
		tmpl.Execute(w, data_obj)
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

// func handleArtists(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		url, _, errUrl := generateUrl(r.URL.Path, "artists")
// 		if errUrl == "not found" {
// 			handleErrorPage(w, r, NotFoundError)
// 			return
// 		}

// 		tmpl, err := template.ParseFiles(
// 			publicUrl+"artists.html",
// 			publicUrl+"header.html",
// 			publicUrl+"artists_content.html",
// 		)
// 		if err != nil {
// 			handleErrorPage(w, r, InternalServerError)
// 			return
// 		}
// 		var data_obj_array []ArtistsData

// 		sendGetRequest(url, &data_obj_array)

// 		tmpl.Execute(w, data_obj_array)

// 	} else {
// 		handleErrorPage(w, r, MethodNotAllowedError)
// 	}
// }

func handleArtist(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "artist")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}

		tmpl, err := template.ParseFiles(
			publicUrl+"artist.html",
			publicUrl+"artist_info.html",
			publicUrl+"artist_dates.html",
			publicUrl+"artist_locations.html",
			publicUrl+"artist_relation.html",
		)
		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}

		var data_obj ArtistsData
		sendGetRequest(url+"/"+id, &data_obj)

		var date_data_obj DatesDataLevel2
		sendGetRequest(data_obj.ConcertDates, &date_data_obj)

		var location_data_obj LocationsDataLevel2
		sendGetRequest(data_obj.Locations, &location_data_obj)
		sendGetRequest(data_obj.ConcertDates, &date_data_obj)

		var relation_data_obj RelationsDataLevel2
		sendGetRequest(data_obj.Relations, &relation_data_obj)

		templateData := struct {
			ArtistInfo      ArtistsData
			ArtistDates     DatesDataLevel2
			ArtistLocations LocationsDataLevel2
			Relation        RelationsDataLevel2
		}{
			ArtistInfo:      data_obj,
			ArtistDates:     date_data_obj,
			ArtistLocations: location_data_obj,
			Relation:        relation_data_obj,
		}

		fmt.Printf("%+v\n", templateData)

		tmpl.Execute(w, templateData)

	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleLocations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "locations")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println(url, id)

		var data_obj LocationsDataLevel1
		sendGetRequest(url, &data_obj)

		fmt.Println(data_obj)
		// fmt.Println("\n________________________________________")

		// for _, data := range data_obj {
		// 	fmt.Println("Id: ", data.Id)
		// 	fmt.Println("Image: ", data.Image)
		// 	fmt.Println("Name: ", data.Name)
		// 	fmt.Println("CreationDate: ", data.CreationDate)
		// 	fmt.Println("FirstAlbum: ", data.FirstAlbum)
		// 	fmt.Println("Locations: ", data.Locations)
		// 	fmt.Println("ConcertDates: ", data.ConcertDates)
		// 	fmt.Println("Relations: ", data.Relations)
		// 	for _, value := range data.Members {
		// 		fmt.Println("Members: ")
		// 		fmt.Print(value)
		// 	}
		// 	fmt.Println("\n________________________________________\n")
		// }

		tmpl, err := template.ParseFiles(publicUrl + "locations.html")

		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}
		tmpl.Execute(w, data_obj)
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleLocation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "location")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println("url is: ", url, " id is: ", id)
		tmpl, err := template.ParseFiles(publicUrl + "locations.html")
		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}
		var data_obj_array []LocationsDataLevel2

		var data_obj LocationsDataLevel2
		sendGetRequest(url+"/"+id, &data_obj)
		fmt.Println("is here 2")
		fmt.Println(data_obj)

		data_obj_array = []LocationsDataLevel2{data_obj}

		fmt.Println("is here 2.5")
		fmt.Println(data_obj_array)

		tmpl.Execute(w, data_obj_array)

	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleDates(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "dates")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println(url, id)

		if id == "" {
			var data_obj []DatesDataLevel1
			sendGetRequest(url, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "dates.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		} else {
			var data_obj DatesDataLevel1
			sendGetRequest(url+"/"+id, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "dates.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		}
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleDate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "date")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println("url is: ", url, " id is: ", id)
		tmpl, err := template.ParseFiles(publicUrl + "dates.html")
		if err != nil {
			handleErrorPage(w, r, InternalServerError)
			return
		}
		var data_obj_array []DatesDataLevel1

		var data_obj DatesDataLevel1
		sendGetRequest(url+"/"+id, &data_obj)
		fmt.Println("is here 2")
		fmt.Println(data_obj)

		data_obj_array = []DatesDataLevel1{data_obj}

		fmt.Println("is here 2.5")
		fmt.Println(data_obj_array)

		tmpl.Execute(w, data_obj_array)

	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleRelations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "relations")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println(url, id)

		if id == "" {
			var data_obj []RelationsDataLevel1
			sendGetRequest(url, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "relations.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		} else {
			var data_obj RelationsDataLevel1
			sendGetRequest(url+"/"+id, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "relations.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		}
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleRelation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		url, id, errUrl := generateUrl(r.URL.Path, "relation")
		if errUrl == "not found" {
			handleErrorPage(w, r, NotFoundError)
			return
		}
		fmt.Println(url, id)

		if id == "" {
			var data_obj []RelationsDataLevel1
			sendGetRequest(url, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "relations.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		} else {
			var data_obj RelationsDataLevel1
			sendGetRequest(url+"/"+id, &data_obj)
			fmt.Println(data_obj)

			tmpl, err := template.ParseFiles(publicUrl + "relations.html")

			if err != nil {
				handleErrorPage(w, r, InternalServerError)
				return
			}
			tmpl.Execute(w, data_obj)
		}
	} else {
		handleErrorPage(w, r, MethodNotAllowedError)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/home" {
			// If the URL is not exactly "/", respond with 404
			handleErrorPage(w, r, NotFoundError)
			return
		}

		req, err := http.NewRequest("GET", apiUrls["artists"], nil)
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

		var data_obj []ArtistsData
		jsonErr := json.Unmarshal(body, &data_obj)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		tmpl, err := template.ParseFiles(publicUrl + "home.html")

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
	getApiUrls()
	http.Handle("/static/", http.FileServer(http.Dir("./frontend/public/")))
	http.Handle("/img/", http.FileServer(http.Dir("./frontend/public/")))

	http.HandleFunc("/", handleIndex)

	// http.HandleFunc("/artists", handleArtists)
	http.HandleFunc("/artist/", handleArtist) //for dynamic routes

	http.HandleFunc("/locations", handleLocations)
	http.HandleFunc("/location/", handleLocation)

	http.HandleFunc("/dates", handleDates)
	http.HandleFunc("/date/", handleDate)

	http.HandleFunc("/relations", handleRelations)
	http.HandleFunc("/relation/", handleRelation)

	http.HandleFunc("/home", handleHome)
	// http.HandleFunc("/ascii-web", handleAsciiWeb)
	// Start the server on port 8082
	fmt.Println("Starting server on 0.0.0.0:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		fmt.Println("we are here")
	}
}
