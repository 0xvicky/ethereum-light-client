package fetchprivatekey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

const privateKeyPath = "./nodeKey"

func savePrivateKey(privateKey *ecdsa.PrivateKey) error {
	//convert the privatekey into bytes
	keyBytes := crypto.FromECDSA(privateKey)

	//convert the keybytes into hexstring
	hexKey := hex.EncodeToString(keyBytes)

	return os.WriteFile(privateKeyPath, []byte(hexKey), 0600)
}
func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	//fetch the hex value
	hexKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		// log.Fatalf("Error occured while loading the hex key:%v", err)
		println("Error occured while loading the hex Key, creating new one...")
		return nil, err
	}
	//decoding it to key bytes

	key, err := hex.DecodeString(string(hexKey))
	if err != nil {
		return nil, err
	}
	// Parse the key bytes into an ECDSA private key
	return crypto.ToECDSA(key)
}

func FetchKey() *ecdsa.PrivateKey {
	var privateKey *ecdsa.PrivateKey
	//generate a new private key
	privateKey, err := loadPrivateKey()
	if err != nil {
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Error occured while generating private key :%v", err)
		}

		if err = savePrivateKey(privateKey); err != nil {
			log.Fatalf("Error occured while saving private key :%v", err)
		}

		println("New key saved!!")
	}
	println("Loaded private key")

	return privateKey
}
