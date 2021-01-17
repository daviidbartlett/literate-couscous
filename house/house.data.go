package house

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)


var houseSlice = struct {
	sync.RWMutex
	s []house
}{s:make([]house, 0)}

func init(){
	fmt.Println(("loading houses..."))

	hSlice, err := loadHouseSlice()
if err != nil {
	log.Fatal(err)
}
	houseSlice.s = hSlice
}	

func loadHouseSlice() ([]house, error) {
	fileName := "houses.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)

	var houses []house

	err = json.Unmarshal([]byte(file), &houses)

	if err != nil {
		log.Fatal(err)
	}
	return houses, nil
}

func writeHouseSlice(json []byte) error {
	fileName:= "houses.json"
	return ioutil.WriteFile(fileName, json, 0777)
}



func getNextID()*int{
	highestID := 0
	for _, h := range houseSlice.s {
		if highestID < *h.HouseID{
			highestID = *h.HouseID
		}
	}
	nextID := highestID + 1
	return &nextID
}

func getHouse(houseID int) *house {
	houseSlice.RLock()
	defer houseSlice.RUnlock()

	for _, h := range houseSlice.s {
		if *h.HouseID == houseID {
			return &h
		}
	}
	return nil
}

func findIndex(houseID int) *int {
	houseSlice.Lock()
	defer houseSlice.Unlock()
	
	for i, h := range houseSlice.s {
		if *h.HouseID == houseID {
			return &i
		}
	}
	return nil
}

func updateHouse(houseID int) bool{
	houseSlice.Lock()
	defer houseSlice.Unlock()

	for i, h := range houseSlice.s {
		if *h.HouseID == houseID {
			houseSlice.s[i] = h
			return true
		}
	}
	return false
}

func insertHouse(h house) {
	houseSlice.Lock()
	defer houseSlice.Unlock()
	h.HouseID = getNextID()
	houseSlice.s = append(houseSlice.s, h)
}

func removeHouse(index int) {
	houseSlice.Lock()
	defer houseSlice.Unlock()
	houseSlice.s = append(houseSlice.s[:index], houseSlice.s[index+1:]...)
}