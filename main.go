package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

// LireLignesDepuisUneLigne lit les lignes à partir d'une ligne spécifiée et retourne les n lignes suivantes sous forme de slice de chaînes.
func LireLignesDepuisUneLigne(nomFichier string, ligneDebut int, nombreLignes int) ([]string, error) {
	f, err := os.Open(nomFichier)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	ligneActuelle := 0
	var lignes []string

	for scanner.Scan() {
		ligneActuelle++
		if ligneActuelle >= ligneDebut && len(lignes) < nombreLignes {
			lignes = append(lignes, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lignes, nil
}

// ConstruireAsciiArt génère l'art ASCII pour une chaîne donnée en utilisant la bannière spécifiée.
func ConstruireAsciiArt(texte string, banner string) (string, error) {
	var result strings.Builder
	const lignesParCaractere = 9

	fichierBannieres := map[string]string{
		"shadow":     "shadow.txt",
		"standard":   "standard.txt",
		"thinkertoy": "thinkertoy.txt",
	}

	fichier, ok := fichierBannieres[banner]
	if !ok {
		return "", fmt.Errorf("%s bannière non trouvée", banner)
	}

	lignesTexte := strings.Split(texte, "\n")

	for _, ligne := range lignesTexte {
		caracteres := make([][]string, lignesParCaractere)
		for i := range caracteres {
			caracteres[i] = make([]string, len(ligne))
		}

		for i, c := range ligne {
			nbline := (c-32)*lignesParCaractere + 1
			lignes, err := LireLignesDepuisUneLigne(fichier, int(nbline), lignesParCaractere)
			if err != nil {
				return "", fmt.Errorf("erreur lors de la lecture des lignes pour le caractère %c: %v", c, err)
			}
			for j := 0; j < lignesParCaractere; j++ {
				if j < len(lignes) {
					caracteres[j][i] = lignes[j]
				} else {
					caracteres[j][i] = strings.Repeat(" ", len(ligne)) // Remplir avec des espaces si nécessaire
				}
			}
		}

		for i := 0; i < lignesParCaractere; i++ {
			result.WriteString(strings.Join(caracteres[i], "") + "\n")
		}

		result.WriteString("\n") // Ajouter une ligne vide entre les lignes de texte
	}

	return result.String(), nil
}

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		texte := r.FormValue("text")
		texte = strings.ReplaceAll(texte, "\\n", "\n")
		banner := r.FormValue("banner")

		if texte == "" {
			HandleError(w, r, "Bad request: missing text", http.StatusBadRequest)
			return
		}
		// Générer l'ASCII Art
		asciiArt, err := ConstruireAsciiArt(texte, banner)
		if err != nil {
			HandleError(w, r, "Error generating ASCII art: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Définir les en-têtes pour le téléchargement du fichier
		w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(asciiArt))
	}
}

// HandleError redirects to the error.html page with the provided error message.
func HandleError(w http.ResponseWriter, r *http.Request, errorMessage string, statusCode int) {
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, "Error loading error template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Données à envoyer au template
	data := struct {
		ErrorCode    int
		ErrorMessage string
	}{
		ErrorCode:    statusCode,
		ErrorMessage: errorMessage,
	}

	// Définir le code d'état HTTP
	w.WriteHeader(statusCode)

	tmpl.Execute(w, data)
}

// Gestionnaire pour afficher la page principale avec le résultat de l'art ASCII
func HandleMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		HandleError(w, r, "Page not found", http.StatusNotFound)
		return
	}
	var asciiArt string
	if r.Method == http.MethodPost {
		texte := r.FormValue("text")
		texte = strings.ReplaceAll(texte, "\\n", "\n")
		banner := r.FormValue("banner")

		if texte == "" {
			HandleError(w, r, "Bad request: missing text or banner", http.StatusBadRequest)
			return
		}
		// Check for non-printable characters
		if !isPrintableASCII(texte) {
			HandleError(w, r, "Bad request: non-printable characters detected", http.StatusBadRequest)
			return
		}
		var err error
		asciiArt, err = ConstruireAsciiArt(texte, banner)
		if err != nil {
			HandleError(w, r, "Error generating ASCII art: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		HandleError(w, r, "Error generating ASCII art: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		AsciiArt string
	}{
		AsciiArt: asciiArt,
	}
	tmpl.Execute(w, data)
}

// isPrintableASCII vérifie si une chaîne contient uniquement des caractères ASCII imprimables
func isPrintableASCII(s string) bool {
	for _, r := range s {
		if (r < 32 || r > 126) && r != 13 && r != 10 {
			return false
		}
	}
	return true
}

func main() {
	mux := http.NewServeMux()

	// Serve static files from the "static" directory
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Handle the main page
	mux.HandleFunc("/", HandleMainPage)
	mux.HandleFunc("/download", HandleDownload)
	server := &http.Server{
		Addr:              ":8044",           //adresse du server
		Handler:           mux,               // listes des handlers
		ReadHeaderTimeout: 10 * time.Second,  // temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second,  // temps maximum d'écriture de la réponse
		IdleTimeout:       120 * time.Second, // temps maximum entre deux rêquetes
		MaxHeaderBytes:    1 << 20,           // 1 MB // maximum de bytes que le serveur va lire
	}

	fmt.Println("Serveur démarré sur le port 'http://localhost:8044'")
	server.ListenAndServe()
}
