package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Tibz-Dankan/keep-active/internal/models"
	"github.com/Tibz-Dankan/keep-active/internal/services"
	"github.com/gorilla/mux"
)

func signInAdmin(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	password := user.Password

	if user.Email == "" || user.Password == "" {
		services.AppError("Missing email/password!", 400, w)
		return
	}

	user, err = user.FindByEMail(user.Email)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if user.ID == "" || user.Role != "sys_admin" {
		services.AppError("Invalid email/password!", 400, w)
		return
	}

	passwordMatches, err := user.PasswordMatches(password)
	if err != nil {
		services.AppError(err.Error(), 500, w)
		return
	}

	if !passwordMatches {
		services.AppError("Invalid email/password!", 400, w)
		return
	}

	random := services.NewRandomNumber()
	OPT := random.D6()

	email := services.Email{Recipient: user.Email}

	if err := email.SendOPT(user.Name, OPT, "Sign in OPT"); err != nil {
		log.Println("Error sending opt email:", err)
		services.AppError(err.Error(), 500, w)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "OPT sent to mail!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SignInAdminRoute(router *mux.Router) {
	router.HandleFunc("/signin-admin", signInAdmin).Methods("POST")
}
