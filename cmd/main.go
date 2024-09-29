package main

import (
	"engine"
	"engine/pgnextract"
	"engine/scoutfish"
	"log"
	"net/http"
)

func main() {
	topLevelRoutes := make(map[string]*engine.Route)

	topLevelRoutes["scoutfish"] = scoutfish.New()
	topLevelRoutes["pgnextract"] = pgnextract.New()

	svr := engine.New(topLevelRoutes)
	err := http.ListenAndServe(":8080", svr)
	log.Println(err)
	log.Println("Server is running")
}
