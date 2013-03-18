package channelExample

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/channel"
	"appengine/user"
)

func init() {

	http.HandleFunc("/_ah/channel/connected/", connected)
	http.HandleFunc("/_ah/channel/disconnected/", disconnected)
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

// func dataStore() {
// 	item := &memcache.Item{
// 		Key:   "lyric",
// 		Value: []byte("Oh, give me a home"),
// 	}
// 	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
// 		c.Infof("item with key %q already exists", item.Key)
// 	} else if err != nil {
// 		c.Errorf("error adding item: %v", err)
// 	}

// 	// Change the Value of the item
// 	item.Value = []byte("Where the buffalo roam")
// 	// Set the item, unconditionally
// 	if err := memcache.Set(c, item); err != nil {
// 		c.Errorf("error setting item: %v", err)
// 	}

// 	// Get the item from the memcache
// 	if item, err := memcache.Get(c, "lyric"); err == memcache.ErrCacheMiss {
// 		c.Infof("item not in the cache")
// 	} else if err != nil {
// 		c.Errorf("error getting item: %v", err)
// 	} else {
// 		c.Infof("the lyric is %q", item.Value)
// 	}
// }

func receive(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("g")
	channel.Send(c, key, "go receive!"+time.Now().String())
}

func connected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	c.Infof("connected")
	c.Infof("%s", key)
}

func disconnected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	key := r.FormValue("from")

	c.Infof("disconnected")
	c.Infof("%s", key)
}
