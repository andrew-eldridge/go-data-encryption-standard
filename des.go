package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type KeyRes struct {
	Key string `json:"key"`
}

type EncryptReq struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}

type EncryptRes struct {
	Cipher string `json:"cipher"`
}

type DecryptReq struct {
	Cipher string `json:"cipher"`
	Key    string `json:"key"`
}

type DecryptRes struct {
	Message string `json:"message"`
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

func encrypt(message, key string) string {
	return ""
}

func decrypt(cipher, key string) string {
	return ""
}

// handler for /api/v1/key {get}
func handleGetKey(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Received: /key {get}")

		// get random 64-bit key
		res, err := json.Marshal(&KeyRes{
			Key: getKey(),
		})
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Unable to marshal JSON response.", http.StatusInternalServerError)
		}

		// write to response object
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

// handler for /api/v1/encrypt {post}
func handleEncrypt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Printf("Received: /encrypt {post}\n")

		// parse request body into EncryptReq struct
		decoder := json.NewDecoder(r.Body)
		var req EncryptReq
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request format.", http.StatusBadRequest)
		}

		// get encryption result
		res, err := json.Marshal(&EncryptRes{
			Cipher: encrypt(req.Message, req.Key),
		})

		// write to response object
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}

// handler for /api/v1/decrypt {post}
func handleDecrypt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Printf("Received: /decrypt {post}\n")

		// parse request body into DecryptReq struct
		decoder := json.NewDecoder(r.Body)
		var req DecryptReq
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request format.", http.StatusBadRequest)
		}

		// get decryption result
		res, err := json.Marshal(&DecryptRes{
			Message: encrypt(req.Cipher, req.Key),
		})

		// write to response object
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(res)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
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
