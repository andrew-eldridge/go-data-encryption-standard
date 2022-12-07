package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
		fmt.Println("Received: /encrypt {post}")

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
		fmt.Println("Received: /decrypt {post}")

		// parse request body into DecryptReq struct
		decoder := json.NewDecoder(r.Body)
		var req DecryptReq
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request format.", http.StatusBadRequest)
		}

		// get decryption result
		res, err := json.Marshal(&DecryptRes{
			Message: decrypt(req.Cipher, req.Key),
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
	fmt.Println("Starting DES HTTP server...")

	// attach API endpoint handlers
	ApiBase := "/api/v1"
	http.HandleFunc(ApiBase+"/key", handleGetKey)
	http.HandleFunc(ApiBase+"/encrypt", handleEncrypt)
	http.HandleFunc(ApiBase+"/decrypt", handleDecrypt)

	// listen on port 8080
	fmt.Println("Listening at 127.0.0.1:8080")
	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
