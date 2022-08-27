package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ramdeuter.org/solscraper/db"
	"ramdeuter.org/solscraper/query"
	"ramdeuter.org/solscraper/scraper"
)

type App struct {
	d        db.DB
	handlers map[string]http.HandlerFunc
}

func NewApp(d db.DB, cors bool) App {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}
	techHandler := app.GetTechnologies
	if !cors {
		techHandler = disableCors(techHandler)
	}
	app.handlers["/api/technologies"] = techHandler
	app.handlers["/create"] = app.createAPI
	app.handlers["/save"] = app.saveAPI
	//app.handlers["/apis"]
	app.handlers["/apis/"] = app.getQueryData
	app.handlers["/apis/meta"] = app.getQueryMeta
	app.handlers["/"] = http.FileServer(http.Dir("/webapp")).ServeHTTP
	return app
}

func (a *App) Serve() error {
	for path, handler := range a.handlers {
		http.Handle(path, handler)
	}
	log.Println("Web server is available on port 8080")
	return http.ListenAndServe(":12345", nil)
}

func (a *App) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	technologies, err := a.d.GetTechnologies()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(technologies)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}
}

func (a *App) createAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := query.Query{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	data := scraper.ScrapeData(q)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func (a *App) saveAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	q := query.Query{}
	err := json.NewDecoder(r.Body).Decode(&q)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	a.d.SaveQuery(q)
}

func (a *App) getQueryMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metadata, err := a.d.GetMetadata()
	if err != nil {
		fmt.Printf("%v", err)
	}
	jsonData, _ := json.Marshal(metadata)
	w.Write(jsonData)
}

func (a *App) getQueryData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	paths := strings.Split(r.URL.String(), "/")
	name := paths[len(paths)-1]
	data, err := a.d.GetDataForQuery(name)
	if err != nil {
		fmt.Printf("%v", err)
	}
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

// Needed in order to disable CORS for local development
func disableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h(w, r)
	}
}
