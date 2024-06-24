package main

import (
	"engine"
	"engine/scoutfish"
	"log"
	"net/http"
)

func main() {
	topLevelRoutes := make(map[string]*engine.Route)

	topLevelRoutes["scoutfish"] = scoutfish.New()

	svr := engine.New(topLevelRoutes)
	http.ListenAndServe(":8080", svr)
	log.Println("Server is running")
}

// func ExecCommand(cmdArgs []string) []byte {

// 	output, err := exec.Command(cmdArgs[0], cmdArgs...).Output()
// 	if err != nil {
// 		panic("could not run executable")
// 	}
// 	return output
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hi, this is he very simple version of the parser service!")
// }

// func main() {
// 	mux := $mux{}
// 	http.HandleFunc("/", handler)
// 	fmt.Println("Parser service is running on port 8080")
// 	http.ListenAndServe(":8080", nil)
// }
