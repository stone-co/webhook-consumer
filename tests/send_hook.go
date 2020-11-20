package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/stone-co/webhook-consumer/pkg/common/keys"
	"gopkg.in/square/go-jose.v2"
)

type Data struct {
	Env             string    `json:"env"`
	EventType       string    `json:"event_type"`
	EventHappenedAt time.Time `json:"event_happened_at"`
	TargetData      struct {
		AccountID    string    `json:"account_id"`
		Amount       int       `json:"amount"`
		BalanceAfter int       `json:"balance_after"`
		CreatedAt    time.Time `json:"created_at"`
		Description  string    `json:"description"`
		ID           string    `json:"id"`
		Operation    string    `json:"operation"`
		Status       string    `json:"status"`
		Type         string    `json:"type"`
	} `json:"target_data"`
}

func main() {
	data := Data{}
	data.Env = "sandbox"
	data.EventType = "cash_in_internal_transfer"
	data.TargetData.AccountID = "09c016b2-876a-450a-9f40-316f8e2f8778"

	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("marshaling data: %v", err)
	}

	encryptedBody := EncryptText("partner/mykey.pub", body)
	signedBody := SignText("stone/mykey.pem", encryptedBody)
	log.Println(signedBody)
	SendHook(signedBody)
}

func EncryptText(keyFile string, text []byte) string {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("reading file %s: %v", keyFile, err)
	}

	pub, err := keys.LoadPublicKey(keyBytes)
	if err != nil {
		log.Fatalf("unable to read public key: %v", err)
	}

	alg := jose.KeyAlgorithm("RSA-OAEP-256")
	enc := jose.ContentEncryption("A256GCM")

	crypter, err := jose.NewEncrypter(enc, jose.Recipient{Algorithm: alg, Key: pub}, nil)
	if err != nil {
		log.Fatalf("unable to instantiate encrypter: %v", err)
	}

	obj, err := crypter.Encrypt(text)
	if err != nil {
		log.Fatalf("unable to encrypt: %v", err)
	}

	// msg := obj.FullSerialize()
	msg, err := obj.CompactSerialize()
	if err != nil {
		log.Fatalf("unable to serialize message: %v", err)
	}

	return msg
}

func SignText(keyFile string, text string) string {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("reading file %s: %v", keyFile, err)
	}

	signingKey, err := keys.LoadPrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	alg := jose.SignatureAlgorithm("PS256")
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: alg, Key: signingKey}, nil)
	if err != nil {
		log.Fatalf("unable to make signer: %v", err)
	}

	obj, err := signer.Sign([]byte(text))
	if err != nil {
		log.Fatalf("unable to sign: %v", err)
	}

	// msg := obj.FullSerialize()
	msg, err := obj.CompactSerialize()
	if err != nil {
		log.Fatalf("unable to serialize message: %v", err)
	}

	return msg
}

func SendHook(encryptedBody string) error {
	requestBody, err := json.Marshal(map[string]string{
		"encrypted_body": encryptedBody,
	})
	if err != nil {
		log.Fatalf("unable to marshal encrypted body: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:3000/api/v0/notifications", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-stone-webhook-event-id", "930bbd6d-0c7a-4fe4-8b50-4b82a20cb847")
	req.Header.Set("x-stone-webhook-event-type", "cash_out_internal_transfer")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("unable to send request to api: %v", err)
	}

	defer resp.Body.Close()

	log.Printf("POST Status: %v\n", resp.Status)

	return nil
}
