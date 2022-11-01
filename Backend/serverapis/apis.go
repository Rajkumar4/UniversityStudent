package serverapis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	db "github.com/UnivertsityStudent/Backend/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var log = logrus.New()

type Server struct {
	Client   *mux.Router
	Database *db.Database
}

type WebData struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Contact   string `json:"contact"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type OutRespose struct {
	FirstName string
	LastName  string
	Contact   string
	Message   string
}

func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	log.SetLevel(logrus.DebugLevel)
	method := r.Method
	if method != "POST" {
		log.Errorf("method is not correct:%s", method)
		w.Write([]byte("method is coreect"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Failed to parse form %s", err.Error())
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(fmt.Sprintf("Failed to parse form %s", err.Error())))
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to parse form data: %s", err.Error())
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(err.Error()))
		return
	}
	d := &WebData{}
	err = json.Unmarshal(data, d)
	if err != nil {
		log.Errorf("Failed to unmarshal body of url %s", err.Error())
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Failed to unmarshal body of url %s", err.Error())))
		return
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("failed to encrypt password %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to encrypt password %s", err.Error())))
		return
	}
	strlist := strings.Split(string(hashedpassword), "-")
	hpass := strings.Join(strlist, "")
	log.Debugf("Check password %s", hpass)
	uid := uuid.Must(uuid.NewRandom()).String()
	err = s.Database.SignUp(uid, d.FirstName, d.LastName, d.Email, d.Contact, hpass)
	if err != nil {
		log.Errorf("Failed to signup users: %s", err.Error())
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(fmt.Sprintf("Failed to signup %s", err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Singup Successfull"))
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	log.SetLevel(logrus.DebugLevel)
	err := r.ParseForm()
	if err != nil {
		log.Errorf("Failed to parse login form %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to parse login form %s", err.Error())))
		return
	}
	login := &LoginData{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read requerst body %s", err.Error())
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(fmt.Sprintf("Failed to read requerst body %s", err.Error())))
		return
	}
	err = json.Unmarshal(body, login)
	if err != nil {
		log.Errorf("failed to unmarshl form data %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to unmarshl form data %s", err.Error())))
		return
	}
	mp, err := s.Database.UserLogin(login.Email)
	if err != nil {
		log.Errorf("Failed to login user %s due to %s", login.Email, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to login user %s due to %s", login.Email, err.Error())))
		return
	}
	// log.Debugf("databse out %s",mp)
	if err := bcrypt.CompareHashAndPassword([]byte(mp["password"]), []byte(login.Password)); err != nil {
		log.Errorf("Password is not correct %s", err.Error())
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		w.Write([]byte(fmt.Sprintf("Password is not correct %s", err.Error())))
		return
	}
	jsonresp := &OutRespose{FirstName: mp["firstname"],
		LastName: mp["lastname"],
		Contact:  mp["contact"],
		Message:  "Welcome to University Postal "}
	res,err := json.Marshal(jsonresp)
	if err!=nil{
		log.Errorf("failed to encode json data %s",err.Error())
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(fmt.Sprintf("failed to encode json data %s",err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
