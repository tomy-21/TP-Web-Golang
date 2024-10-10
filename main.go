package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sync"
)

// Étudiant
type Student struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string // "M" pour Masculin, "F" pour Féminin
}

// Classe
type Classe struct {
	Nom             string
	Filiere         string
	Niveau          string
	NombreEtudiants int
	Etudiants       []Student
}

var tmpl = template.Must(template.ParseFiles("promo.html"))

var (
	viewCounter int
	mu          sync.Mutex
	changeTmpl  = template.Must(template.ParseFiles("change.html"))
)

func promoHandler(w http.ResponseWriter, r *http.Request) {
	classe := Classe{
		Nom:             "Promo B1 Informatique",
		Filiere:         "Informatique",
		Niveau:          "Bachelor 1",
		NombreEtudiants: 5,
		Etudiants: []Student{
			{Nom: "Ouadhi", Prenom: "Tomy", Age: 19, Sexe: "M"},
			{Nom: "Thimothé", Prenom: "Claire", Age: 68, Sexe: "F"},
			{Nom: "Amaru", Prenom: "Tom", Age: 14, Sexe: "M"},
			{Nom: "Bourdeau", Prenom: "Romain", Age: 20, Sexe: "M"},
			{Nom: "Rodrigues", Prenom: "Cyril", Age: 23, Sexe: "M"},
		},
	}

	tmpl.Execute(w, classe)
}

func changeHandler(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	viewCounter++
	count := viewCounter
	mu.Unlock()

	parite := "pair"
	message := fmt.Sprintf("Le nombre de vues est pair : %d", count)
	if count%2 != 0 {
		parite = "impair"
		message = fmt.Sprintf("Le nombre de vues est impair : %d", count)
	}

	data := struct {
		Parite  string
		Message string
	}{
		Parite:  parite,
		Message: message,
	}

	changeTmpl.Execute(w, data)
}

// CHALENG 3

var (
	formTmpl    = template.Must(template.ParseFiles("templates/user_form.html"))
	displayTmpl = template.Must(template.ParseFiles("templates/user_display.html"))
	errorTmpl   = template.Must(template.ParseFiles("templates/user_error.html"))
	currentUser User
)

type User struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	formTmpl.Execute(w, nil)
}

func treatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	nom := r.FormValue("nom")
	prenom := r.FormValue("prenom")
	dateNaissance := r.FormValue("date_naissance")
	sexe := r.FormValue("sexe")

	isValid := validateUserData(nom, prenom, sexe)
	if !isValid {
		http.Redirect(w, r, "/user/error", http.StatusSeeOther)
		return
	}

	currentUser = User{
		Nom:           nom,
		Prenom:        prenom,
		DateNaissance: dateNaissance,
		Sexe:          sexe,
	}

	http.Redirect(w, r, "/user/display", http.StatusSeeOther)
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	if currentUser.Nom == "" || currentUser.Prenom == "" || currentUser.DateNaissance == "" || currentUser.Sexe == "" {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	displayTmpl.Execute(w, currentUser)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	errorTmpl.Execute(w, nil)
}

func validateUserData(nom, prenom, sexe string) bool {
	validSexe := sexe == "masculin" || sexe == "féminin" || sexe == "autre"
	validNomPrenom := regexp.MustCompile(`^[A-Za-z]{1,32}$`).MatchString

	return validNomPrenom(nom) && validNomPrenom(prenom) && validSexe
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/promo", promoHandler)

	http.HandleFunc("/change", changeHandler)

	http.HandleFunc("/user/form", formHandler)
	http.HandleFunc("/user/treatment", treatmentHandler)
	http.HandleFunc("/user/display", displayHandler)
	http.HandleFunc("/user/error", errorHandler)

	http.ListenAndServe(":8080", nil)
}
