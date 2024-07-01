package delivery

import "github.com/gorilla/mux"

func UserRoutes(route *mux.Router, userHandler *UserHandler) {
	route.HandleFunc("/login-google", userHandler.GoogleLogin)
	route.HandleFunc("/google-callback", userHandler.CallBackGoogle)
}