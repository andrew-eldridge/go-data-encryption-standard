package main

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
