package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vchatchai/tracking/server/db"
	"github.com/vchatchai/tracking/server/generate"
	"github.com/vchatchai/tracking/server/model"
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
	bookingRefreshHandler := app.RefreshBooking
	containerHandler := app.GetContainner
	containerRefreshHandler := app.RefreshContainer
	loginHandler := app.Login
	loginRefreshHandler := app.RefreshUser
	if !cors {
		bookingHandler = disableCors(bookingHandler)
		containerHandler = disableCors(containerHandler)
		loginHandler = disableCors(loginHandler)
	}

	app.handlers["/refresh/booking"] = bookingRefreshHandler
	app.handlers["/refresh/container"] = containerRefreshHandler
	app.handlers["/refresh/user"] = loginRefreshHandler

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
		port = ":8081"
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

		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}

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
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}

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
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
		}
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

func (a *App) RefreshBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	apiKey := r.Header.Get("x-api-key")

	fmt.Println("APIKEY:", apiKey)

	fmt.Println("Test RefreshBooking")
	data, _ := ioutil.ReadAll(r.Body)
	// println(string(data))

	var bookings []*model.Booking

	// err := json.NewDecoder(r.Body).Decode(&bookings)
	err := json.Unmarshal([]byte(data), &bookings)

	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}
	// fmt.Println("bookings len", len(bookings))
	// for _, booking := range bookings {
	// 	// 	return
	// 	fmt.Println("booking:", booking)
	// }

	err = a.d.RefreshBooking(bookings)

	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	// fmt.Fprintln(io.)
	// json.NewDecoder(r.Body).Decode(book)

}

func (a *App) RefreshContainer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	apiKey := r.Header.Get("x-api-key")

	fmt.Println("APIKEY:", apiKey)
	data, _ := ioutil.ReadAll(r.Body)
	// println(string(data))

	var containers []*model.LadenContainer

	// err := json.NewDecoder(r.Body).Decode(&containers)
	err := json.Unmarshal([]byte(data), &containers)

	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.d.RefreshContainer(containers)

	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("LadenContainer refreshed")

}

func (a *App) RefreshUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	apiKey := r.Header.Get("x-api-key")

	fmt.Println("APIKEY:", apiKey)
	data, _ := ioutil.ReadAll(r.Body)
	println(string(data))

	var users []*model.User

	// err := json.NewDecoder(r.Body).Decode(&users)
	err := json.Unmarshal([]byte(data), &users)
	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("user len", len(users))
	// for _, user := range users {
	// 	// 	return
	// 	fmt.Println("User:", user)
	// }

	err = a.d.RefreshUser(users)

	if err != nil {
		sendErr(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("user refreshed")

	// fmt.Fprintln(io.)
	// json.NewDecoder(r.Body).Decode(book)

}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
	log.Println(string(resp))

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
