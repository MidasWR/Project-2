package api

import (
	"PROJECT-2/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Server struct {
	config *config.Config
	router *mux.Router
}

var SignedKey = []byte("secret")

func (s *Server) NewRouter() (error, error) {
	s.router.Handle("/", IsAuthorized(s.HomePage))
	s.router.HandleFunc("/login", s.LogPage()).Methods("POST")
	s.router.HandleFunc("/registration", s.RegPage()).Methods("POST")
	log.Println("server opening access")
	return http.ListenAndServe(":8060", s.router), s.config.DB.Close()
}

func (s *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Home page")
	if err != nil {
		log.Println("error writing")
		panic(err)
	}
}

func (s *Server) LogPage() http.HandlerFunc {
	var credentials struct {
		Login    string
		Password string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		credentials.Login = r.URL.Query().Get("login")
		credentials.Password = r.URL.Query().Get("password")
		log.Println(LoginUsers(s.config.DB, credentials.Login, credentials.Password))
		log.Println("Access login in!")
	}
}
func (s *Server) RegPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		login := r.URL.Query().Get("login")
		password := r.URL.Query().Get("password")
		if err := RegUser(*s.config, username, login, password); err != nil {
			log.Println("error reg in")
		}
		if _, err := io.WriteString(w, "RegPage"); err != nil {
			log.Println("error writing")
		}
	}
}
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return SignedKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
func NewAPI() *Server {
	return &Server{
		config: config.NewConfig(),
		router: mux.NewRouter(),
	}
}
