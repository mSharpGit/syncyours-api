// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//<username>:<pw>@tcp(<HOST>:<port>)/<dbname>
func (a *App) Initialize(user, password, host, port, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.DB.SetConnMaxLifetime(1000 * time.Millisecond)
	//defer a.DB.Close()

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/authuser", a.authUser).Methods("POST")
	a.Router.HandleFunc("/resetuser", a.resetUser).Methods("POST")
	a.Router.HandleFunc("/resetuser/{id:[0-9]+}/{code}", a.updateUserPass).Methods("POST")
	a.Router.HandleFunc("/valuser/{id:[0-9]+}/{code}", a.valUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

func (a *App) updateUserPass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	code := vars["code"]
	/* if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	} */
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	if err := u.updateUserPass(a.DB, id, code); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "Wrong Password Update Request")
		default:

			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)

}

func (a *App) resetUser(w http.ResponseWriter, r *http.Request) {

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.resetUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) authUser(w http.ResponseWriter, r *http.Request) {

	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.authUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			var ErrHashTooShort = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
			if err.Error() == ErrHashTooShort.Error() {
				respondWithError(w, http.StatusNotFound, "Password is wrong")
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
			}

		}
		return
	}
	respondWithJSON(w, http.StatusCreated, u)

	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
}

func (a *App) valUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	code := vars["code"]
	//log.Println("id:", id)
	//log.Println("code:", code)

	/* if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	} */
	var u user
	u.ID = id
	u.Verifycode = code
	/* decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close() */
	if err := u.verifyUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:

			respondWithError(w, http.StatusNotFound, "Wrong Validation Request")
		default:

			respondWithError(w, http.StatusInternalServerError, err.Error())

		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)

}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	u := user{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Println("error: ", message)
	respondWithJSON(w, code, map[string]string{"error": message})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	response, _ := json.Marshal(payload)
	log.Println("response: ", payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	users, err := getUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var u user
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id
	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	u := user{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
