// main.go
package main

import (
	"log"
	"net/http"

	"github.com/clover-network/ethereum-spv-backend/app/controller"
	"github.com/gorilla/mux"
)

// func verifyTransaction(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	txID := vars["id"]

// 	_, merklePath, err := controller.VerifyTransaction(txID)
// 	if err != nil {
// 		merklePath = nil
// 	}

// 	m := make(map[string]string)
// 	for k, v := range merklePath {
// 		value := fmt.Sprintf("0x%x", v)
// 		m["0x"+k] = value
// 	}

// 	json.NewEncoder(w).Encode(m)
// }

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/merklepath/{id}", verifyTransaction)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	// handleRequests()
	// app.ListenToEvents()
	// controller.VerifyTransaction("0xaaabc2afc8b1efe6a386f4ef9826df9c5122c29aeb72b6bcfc7f0658e5357548")
	controller.GetDerive("0xaaabc2afc8b1efe6a386f4ef9826df9c5122c29aeb72b6bcfc7f0658e5357548")
}
