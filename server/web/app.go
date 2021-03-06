package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tracking/db"
	"tracking/generate"
	"tracking/model"
)

type App struct {
	d        db.DB
	handlers map[string]http.HandlerFunc
}

var Environment model.Environment

func NewApp(d db.DB, cors bool) App {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}

	bookingHandler := app.GetBooking
	containerHandler := app.GetContainner
	loginHandler := app.Login
	if !cors {
		bookingHandler = disableCors(bookingHandler)
		containerHandler = disableCors(containerHandler)
		loginHandler = disableCors(loginHandler)
	}

	app.handlers["/booking/"] = bookingHandler
	app.handlers["/container/"] = containerHandler
	app.handlers["/login"] = loginHandler
	app.handlers["/"] = http.FileServer(generate.Assets).ServeHTTP
	return app
}

func (a *App) Serve() error {
	for path, handler := range a.handlers {
		http.Handle(path, handler)
	}

	port := fmt.Sprintf(":%d", Environment.HttpPort)
	if Environment.HttpPort == 0 {
		port = ":8080"
	}
	log.Println("Web server is available on port :", port)
	return http.ListenAndServe(port, nil)
}

func (a *App) GetContainner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// keys, ok := r.URL.Query()["bookingNumber"]
	// keys, ok := r.URL.Query()["containerNumber"]
	bookingNumbers, bookingOk := r.URL.Query()["bookingNumber"]
	containerNumbers, containerOk := r.URL.Query()["containerNumber"]

	if bookingOk || len(bookingNumbers) >= 1 {

		key := bookingNumbers[0]

		log.Println("GetContainner 'bookingNumber' is: " + string(key))

		booking, err := a.d.GetContainerByBookingNumber(key)

		err = json.NewEncoder(w).Encode(booking)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if containerOk || len(containerNumbers[0]) >= 1 {

		key := containerNumbers[0]

		log.Println("GetContainner 'containerNumber' is: " + string(key))

		booking, err := a.d.GetContainerContainerNumber(key)

		err = json.NewEncoder(w).Encode(booking)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
}

func (a *App) GetBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookingNumbers, bookingOk := r.URL.Query()["bookingNumber"]
	containerNumbers, containerOk := r.URL.Query()["containerNumber"]

	if bookingOk || len(bookingNumbers) >= 1 {
		key := bookingNumbers[0]

		log.Println("GetBooking 'bookingNumber' is: " + string(key))

		booking, err := a.d.GetBookingByBookingNumber(key)
		err = json.NewEncoder(w).Encode(booking)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if containerOk || len(containerNumbers[0]) >= 1 {
		key := containerNumbers[0]

		log.Println("GetBooking 'containerNumber' is: " + string(key))

		booking, err := a.d.GetBookingByContainerNumber(key)
		err = json.NewEncoder(w).Encode(booking)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, ok := r.URL.Query()["username"]
	passwords, ok := r.URL.Query()["password"]
	if !ok || len(users[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	user := users[0]
	password := passwords[0]

	result, err := a.d.Login(user, password)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}

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
