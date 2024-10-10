package main

import (
	"fmt"
	"html/template"
	"net/http"
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

func changeHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	// Sert les fichiers statiques (images et CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route pour afficher les promotions
	http.HandleFunc("/promo", promoHandler)

	// Route pour afficher le challenge 2 - Compteur de vues
	http.HandleFunc("/change", changeHandler)

	// Démarrer le serveur
	http.ListenAndServe(":8080", nil)
}
