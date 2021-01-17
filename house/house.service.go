package house

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Welcome handler triggered on /api path
func Welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w,"Hello and welcome to my Harry Potter API")
}

// HandleHouses triggered on /api/houses for GET & POST 
func HandleHouses(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		
	// GET
	case http.MethodGet:
		housesJSON, err := json.Marshal(houseSlice.s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(housesJSON)
	// POST
	case http.MethodPost:
		bodyBytes, err := ioutil.ReadAll(r.Body)

		var newHouse house

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newHouse)
		
		if err != nil || newHouse.HouseID != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		insertHouse(newHouse)

		housesJSON, err := json.MarshalIndent(houseSlice.s,"", "  ")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = writeHouseSlice(housesJSON)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		newHouseJSON, err := json.Marshal(newHouse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(newHouseJSON)
	default: w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
// HandleHouse triggered on /api/houses/:house_id for GET, PATCH & DELETE
func HandleHouse(w http.ResponseWriter, r *http.Request) {

	pathComponents:= strings.Split(r.URL.Path,"/")
	id, err := strconv.Atoi(pathComponents[len(pathComponents)-1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
 
	switch r.Method {
	// GET
	case http.MethodGet:
		 h := getHouse(id) 
		 fmt.Println(h)
		if h == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	houseJSON, err := json.Marshal(h)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(houseJSON)
	// PUT
	case http.MethodPut:
		body, _ := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var h house
		err := json.Unmarshal(body,&h)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if h.HouseID == nil {
			fmt.Println(h)
			h.HouseID = &id
		}
		


		if ok:= updateHouse(id); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		
		jsonHouses, err := json.MarshalIndent(houseSlice.s,"", "  ")
		
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = writeHouseSlice(jsonHouses)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonHouses)
	// DELETE
	case http.MethodDelete:
	
	index := findIndex(id)
	
	if index == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	removeHouse(*index)

	housesJSON, err := json.MarshalIndent(houseSlice.s, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = writeHouseSlice(housesJSON)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	default: w.WriteHeader(http.StatusMethodNotAllowed)
}
}