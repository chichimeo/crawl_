package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chichimeo/crawl/api"
	"github.com/chichimeo/crawl/crawl"
	"github.com/chichimeo/crawl/malware"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type Session struct {
	session *mgo.Session
}

type Config struct {
	Hosts      string
	Database   string
	UserName   string
	Password   string
	Collection string
	Server     string
}

func NewSession(conf Config) (*Session, error) {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{conf.Hosts},
		Database: conf.Database,
		Username: conf.UserName,
		Password: conf.Password,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return &Session{session}, err
	}
	session.SetMode(mgo.Monotonic, true)

	return &Session{session}, nil
}

func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

func main() {
	var conf Config
	_, err := toml.DecodeFile("config_example.toml", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := NewSession(conf)
	defer session.Close()

	collectionRepository := malware.NewMongoRepository(session.Copy().DB(conf.Database))

	switch os.Args[1] {
	case "crawl":
		malwareHandler := crawl.MalwareHandler{
			Repository: collectionRepository,
		}
		malwareHandler.Crawl()
	case "resful":
		malwareHandler := api.MalwareHandler{
			Repository: collectionRepository,
		}

		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.HandleFunc("/api/malware/", malwareHandler.GetMalwares).Methods("GET")
		myRouter.HandleFunc("/api/malware/{id}", malwareHandler.GetByID).Methods("GET")
		myRouter.HandleFunc("/api/malware/md5/{hash}", malwareHandler.GetOneByMd5).Methods("GET")
		myRouter.HandleFunc("/api/malware/sha1/{hash}", malwareHandler.GetOneBySha1).Methods("GET")
		myRouter.HandleFunc("/api/malware/sha256/{hash}", malwareHandler.GetOneBySha256).Methods("GET")
		myRouter.HandleFunc("/api/malware/", malwareHandler.CreateData).Methods("POST")
		myRouter.HandleFunc("/api/malware/{hash}", malwareHandler.UpdateData).Methods("PATCH")
		myRouter.HandleFunc("/api/malware/{hash}", malwareHandler.UpdateData).Methods("PUT")
		myRouter.HandleFunc("/api/malware/{hash}", malwareHandler.DeleteData).Methods("DELETE")
		err = http.ListenAndServe(conf.Server, myRouter)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
