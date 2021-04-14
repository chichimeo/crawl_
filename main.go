package main

import (
	"fmt"
	"net/http"

	"Crawl/malware"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Print("Enter case: ")
	text := ""
	fmt.Scanln(&text)
	switch text {
	case "crawl":
		fmt.Println(text)
		collection, close := malware.NewCollection()
		defer close()
		malware.Crawl(collection)
		return
	case "resful":
		myRouter := mux.NewRouter().StrictSlash(true)
		myRouter.HandleFunc("/api/malware/", malware.GetMalwares).Methods("GET")
		myRouter.HandleFunc("/api/malware/{id}", malware.GetByID).Methods("GET")
		myRouter.HandleFunc("/api/malware/md5/{hash}", malware.GetOneByMd5).Methods("GET")
		myRouter.HandleFunc("/api/malware/sha1/{hash}", malware.GetOneBySha1).Methods("GET")
		myRouter.HandleFunc("/api/malware/sha256/{hash}", malware.GetOneBySha256).Methods("GET")
		myRouter.HandleFunc("/api/malware/", malware.CreateData).Methods("POST")
		myRouter.HandleFunc("/api/malware/{hash}", malware.UpdateDataPatch).Methods("PATCH")
		myRouter.HandleFunc("/api/malware/{hash}", malware.UpdateDataPut).Methods("PUT")
		myRouter.HandleFunc("/api/malware/{hash}", malware.DeleteData).Methods("DELETE")

		err := http.ListenAndServe("127.0.0.1:8081", myRouter)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
