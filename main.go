package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type house struct {
	HouseID *int `json:"house_id,omitempty"`
	Name string `json:"house_name,omitempty"`
	Founder string `json:"founder,omitempty"`
	Animal string `json:"animal,omitempty"`
}

var houses []house

func init(){
housesJSON, _ := ioutil.ReadFile("houses.json")
	err := json.Unmarshal([]byte(housesJSON), &houses)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID ()*int{
	highestID := 0
	for _, house := range houses {
		if highestID < *house.HouseID{
			highestID = *house.HouseID
		}
	}
	nextID := highestID + 1
	return &nextID
}


func welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w,"Hello and welcome to my Harry Potter API")
}

func handleHouses(w http.ResponseWriter, r *http.Request) {
	
	switch r.Method {
		
	// GET
	case http.MethodGet:
		housesJSON, err := json.Marshal(houses)
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
		
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newHouse.HouseID != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		newHouse.HouseID = getNextID()
		houses = append(houses,newHouse)

		housesJSON, err := json.MarshalIndent(houses,"", "  ")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ioutil.WriteFile("houses.json", housesJSON,0777)

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

func handleHouse(w http.ResponseWriter, r *http.Request) {

	pathComponents:= strings.Split(r.URL.Path,"/")
	id, err := strconv.Atoi(pathComponents[len(pathComponents)-1])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
 
	switch r.Method {
	// GET
	case http.MethodGet:
		var retrievedHouse house

		for _, house := range houses {
			if *house.HouseID == id {
				retrievedHouse = house
				break
			}
		}

		if retrievedHouse == (house{}) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	houseJSON, err := json.Marshal(retrievedHouse)
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
		


		ok:= false
		for i, house := range houses {
			if *house.HouseID == id {
				houses[i] = h
				ok = true
				break
			}
		}
		// check that it found a house to update
		if ok != true {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonHouses, err := json.MarshalIndent(houses,"", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ioutil.WriteFile("houses.json", jsonHouses, 0777)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonHouses)
	// DELETE
	case http.MethodDelete:
		var iToRemove *int
	for i, house := range houses {
		if *house.HouseID == id {
			iToRemove = &i
			break
		}
	}
	if iToRemove == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	houses = append(houses[:*iToRemove], houses[*iToRemove+1:]...)
	fmt.Println(houses)

	housesJSON, err := json.MarshalIndent(houses, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = ioutil.WriteFile("houses.json", housesJSON, 0777)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	default: w.WriteHeader(http.StatusMethodNotAllowed)
}
	
}


			   

func main() {

	
http.HandleFunc("/api", welcome)
http.HandleFunc("/api/houses", handleHouses)
http.HandleFunc("/api/houses/", handleHouse)

http.ListenAndServe(":9090",nil)

}