package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type KeyRes struct {
	key string
}

func randBits(nBits int) string {
	base2 := []rune("01")
	bits := make([]rune, nBits)
	for i := range bits {
		bits[i] = base2[rand.Intn(2)]
	}
	return string(bits)
}

func getKey() string {
	return randBits(64)
}

// handler for /api/v1/key {get}
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		key := getKey()
		fmt.Println(key)
		res, err := json.Marshal(KeyRes{
			key: key,
		})
		fmt.Println(res)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Unable to marshal JSON response.", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		//_, err = w.Write(res)
		json.NewEncoder(w).Encode(res)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
	fmt.Printf("/key {get}\n")
}

// handler for /api/v1/encrypt {post}
func handleEncrypt(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/encrypt {post}\n")
}

// handler for /api/v1/decrypt {post}
func handleDecrypt(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("/decrypt {post}\n")
}

func main() {
	// attach API endpoint handlers
	ApiBase := "/api/v1"
	http.HandleFunc(ApiBase+"/key", handleGetKey)
	http.HandleFunc(ApiBase+"/encrypt", handleEncrypt)
	http.HandleFunc(ApiBase+"/decrypt", handleDecrypt)

	// listen on port 8080
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
