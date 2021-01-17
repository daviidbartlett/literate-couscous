package main

import (
	"net/http"
	"wizards/house"
)


func main() {
http.HandleFunc("/api", house.Welcome)
http.HandleFunc("/api/houses", house.HandleHouses)
http.HandleFunc("/api/houses/", house.HandleHouse)

http.ListenAndServe(":9090",nil)

}