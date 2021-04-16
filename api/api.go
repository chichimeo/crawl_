package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/chichimeo/Crawl/malware"

	"github.com/gorilla/mux"
)

type MalwareHandler struct {
	Repository malware.Repository
}

const (
	AddedMessage   = "Added %s"
	UpdatedMessage = "Updated %s"
	DeletedMessage = "Deleted %s"
)

func (re *MalwareHandler) CreateData(w http.ResponseWriter, r *http.Request) {
	var newData malware.Malware
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	json.Unmarshal(reqBody, &newData)

	err = re.Repository.Insert(newData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf(AddedMessage, newData.Md5),
	})
}
func (re *MalwareHandler) UpdateData(w http.ResponseWriter, r *http.Request) {
	var newData malware.Malware
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	json.Unmarshal(reqBody, &newData)
	hash := mux.Vars(r)["hash"]
	err = re.Repository.Update(hash, newData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf(UpdatedMessage, hash),
	})
}

func (re *MalwareHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	err := re.Repository.Delete(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf(AddedMessage, hash),
	})
}

func (re *MalwareHandler) GetOneByMd5(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	data, err := re.Repository.FindByMd5(hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)

}

func (re *MalwareHandler) GetOneBySha1(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	data, err := re.Repository.FindBySha1(hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (re *MalwareHandler) GetOneBySha256(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	data, err := re.Repository.FindBySha256(hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (re *MalwareHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	data, err := re.Repository.FindByID(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (re *MalwareHandler) GetMalwares(w http.ResponseWriter, r *http.Request) {
	skip := 0

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	skip = (page - 1) * limit
	if err != nil || page == 0 {
		skip = 0
	}
	malwareList, err := re.Repository.List(skip, limit)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(malwareList)

}
