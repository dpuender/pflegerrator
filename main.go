package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pflegerrator/service"
	"pflegerrator/structs"
	"time"
)

//go:embed views/*
var views embed.FS

var t = template.Must(template.ParseFS(views, "views/*"))

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "form.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/submit-form", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {

			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Fehler beim Verarbeiten des Formulars", http.StatusBadRequest)
				return
			}

			personData := structs.Person{
				LastName:  r.FormValue("nachname"),
				BirthDate: r.FormValue("geburtsdatum"),
				Sex:       r.FormValue("geschlecht"),
			}

			rvNr := service.RvGenerator(personData)

			personData.RvNr = rvNr

			if err := t.ExecuteTemplate(w, "response.html", personData); err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
		}

	})

	router.HandleFunc("/card", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "card.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	router.HandleFunc("/randomuser", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("https://randomuser.me/api/")

		if err != nil {
			http.Error(w, "Failed to fetch random user", http.StatusInternalServerError)
			log.Printf("HTTP request error: %v", err)
			return
		}
		defer resp.Body.Close()

		data := structs.Results{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			http.Error(w, "Failed to decode API response", http.StatusInternalServerError)
			log.Printf("JSON decode error: %v", err)
			return
		}

		log.Printf("Fetched random user: %+v", data)

		pflegeperson := mapUserToPerson(data)

		if err := t.ExecuteTemplate(w, "randomuser.html", pflegeperson); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Template execution error: %v", err)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server running on port 8080")
	server.ListenAndServe()

}

func mapUserToPerson(data structs.Results) structs.Person {
	var pflegeperson structs.Person
	pflegeperson.FirstName = data.Results[0].Name.First
	pflegeperson.LastName = data.Results[0].Name.Last

	input := data.Results[0].DOB.Date

	parsedTime, err := time.Parse(time.RFC3339, input)
	if err != nil {
		fmt.Println("Fehler beim Parsen des Datums:", err)
	}

	pflegeperson.BirthDate = parsedTime.Format("02.01.2006")

	if data.Results[0].Gender == "male" {
		pflegeperson.Sex = "männlich"
	} else {
		pflegeperson.Sex = "weiblich"
	}

	pflegeperson.Picture = data.Results[0].Picture
	pflegeperson.RvNr = service.RvGenerator(pflegeperson)
	return pflegeperson
}
