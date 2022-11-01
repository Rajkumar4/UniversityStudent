package backend

import (
	"net/http"

	db "github.com/UnivertsityStudent/Backend/database"
	api "github.com/UnivertsityStudent/Backend/serverapis"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitServer() error {
	log.SetLevel(logrus.DebugLevel)
	r := mux.NewRouter()
	d, err := db.Init_Database()
	if err != nil {
		log.Errorf("Failed to init database: %s", err.Error())
		return err
	}
	svr := &api.Server{Client: r,
		Database: d}
	r.HandleFunc("/signup", svr.Signup).Methods("POST")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to univertsity"))
	})
	r.HandleFunc("/login",svr.Login).Methods("POST")
	log.Infof("Server start at 8000")
	http.ListenAndServe(":8000", r)
	defer svr.Database.DB.Close()
	return nil
}
