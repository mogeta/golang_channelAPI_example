package tictactoe

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/channel"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", main)
	http.HandleFunc("/receive", receive)
}

var mainTemplate = template.Must(template.ParseFiles("channel/main.html"))

func main(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c) // assumes 'login: required' set in app.yaml
	key := r.FormValue("gamekey")
	tok, err := channel.Create(c, u.ID+key)
	if err != nil {
		http.Error(w, "Couldn't create Channel", http.StatusInternalServerError)
		c.Errorf("channel.Create: %v", err)
		return
	}

	err = mainTemplate.Execute(w, map[string]string{
		"token":    tok,
		"me":       u.ID,
		"game_key": key,
	})
	if err != nil {
		c.Errorf("mainTemplate: %v", err)
	}
}

func receive(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("g")
	channel.Send(c, key, "go側のreceive!")
}
