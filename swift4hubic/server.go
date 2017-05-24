package swift4hubic

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var templates = template.Must(template.ParseGlob("swift4hubic/templates/*.html"))

type appContext struct {
	config *Configuration
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		"Test",
	}
	renderTemplate(w, "index", data)
}

func (cxt *appContext) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var randomState = "1234567890"

	hubicApp := cxt.config.HubicApplications[0]
	config := NewOAuthConfig(hubicApp)
	url := config.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (cxt *appContext) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var randomState = "1234567890"

	hubicApp := cxt.config.HubicApplications[0]
	config := NewOAuthConfig(hubicApp)

	state := r.FormValue("state")
	if state != randomState {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", randomState, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	//ctx := context.Background()
	// TODO Store the token
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		Log(LogError, err.Error())
	}

	client := config.Client(r.Context(), token)
	response, err := client.Get("https://api.hubic.com/1.0/account/usage")
	if err != nil {
		Log(LogError, err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "Content: %s\n", contents)
}

func (cxt *appContext) AuthHandler(w http.ResponseWriter, r *http.Request) {
	hubicApp := cxt.config.HubicApplications[0]

	token, err := hubicApp.getToken()
	if err != nil {
		//TODO add user to context and redirect to handleLogin
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := NewOAuthConfig(hubicApp).Client(context.Background(), token)

	swiftToken, err := NewSwiftToken(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-Storage-Url", swiftToken.StorageURL)
	w.Header().Set("X-Auth-Token", swiftToken.AuthToken)
	w.WriteHeader(http.StatusNoContent) // 204
}

func NewServer(config *Configuration) {
	appCxt := &appContext{
		config: config,
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/register", appCxt.RegisterHandler)
	http.HandleFunc("/callback/", appCxt.CallbackHandler)
	http.HandleFunc("/auth/v1.0/", appCxt.AuthHandler)

	http.ListenAndServe(":"+config.Port, nil)
}
