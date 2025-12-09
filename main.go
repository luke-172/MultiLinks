package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Link struct {
	Titulo string
	URL    string
}

type Perfil struct {
	Nome       string
	Bio        string
	Foto       string
	Background string
	Links      []Link
}

var (
	templates map[string]*template.Template
	perfil    Perfil
)

func main() {
	templates = make(map[string]*template.Template)
	carregarTemplates()

	// Serve arquivos estáticos (CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Rotas
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/painel", painelPage)
	http.HandleFunc("/salvar", salvarHandler)
	http.HandleFunc("/linktree", linktreePage)

	fmt.Println("Servidor em http://localhost:8080")
	fmt.Println("Usuário: admin / Senha: 123")
	
	http.ListenAndServe(":8080", nil)
}

func carregarTemplates() {
	templates["login"] = template.Must(template.ParseFiles("templates/login.html"))
	templates["painel"] = template.Must(template.ParseFiles("templates/painel.html"))
	templates["linktree"] = template.Must(template.ParseFiles("templates/linktree.html"))
}

func logado(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "ok" {
		return false
	}
	return true
}

// CONTROLADORES
func loginPage(w http.ResponseWriter, r *http.Request) {
	templates["login"].Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	usuario := r.FormValue("usuario")
	senha := r.FormValue("senha")

	if usuario == "admin" && senha == "123" {
		cookie := &http.Cookie{
			Name:  "session",
			Value: "ok",
			Path:  "/",
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/painel", http.StatusSeeOther)
	} else {
		dados := map[string]string{"Erro": "Usuário ou senha incorretos!"}
		templates["login"].Execute(w, dados)
	}
}

func painelPage(w http.ResponseWriter, r *http.Request) {
	if !logado(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates["painel"].Execute(w, perfil)
}

func salvarHandler(w http.ResponseWriter, r *http.Request) {
	if !logado(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Redirect(w, r, "/painel", http.StatusSeeOther)
		return
	}

	perfil.Nome = r.FormValue("nome")
	perfil.Bio = r.FormValue("bio")
	perfil.Foto = r.FormValue("foto")
	bg := r.FormValue("background")
	if bg == "" {
		bg = "linear-gradient(135deg, #667eea, #764ba2)"
	}
	perfil.Background = bg

	var links []Link

	for i := 1; i <= 5; i++ {
		titulo := r.FormValue("titulo" + strconv.Itoa(i))
		url := r.FormValue("url" + strconv.Itoa(i))

		if titulo != "" && url != "" {
			if len(url) >= 4 && url[:4] != "http" {
				url = "https://" + url
			}

			links = append(links, Link{
				Titulo: titulo,
				URL:    url,
			})
		}
	}

	perfil.Links = links
	http.Redirect(w, r, "/painel", http.StatusSeeOther)
}

func linktreePage(w http.ResponseWriter, r *http.Request) {
	templates["linktree"].Execute(w, perfil)
}