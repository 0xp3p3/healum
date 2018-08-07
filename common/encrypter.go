package common

import (
	"github.com/mailgun/lemma/secret"
	"encoding/base64"
	"errors"
)

var (
	DefaultKey = "dGhpc2lzYXNlY3JldGtleXBsZWFzZWRvbm90c2hhcmU="
)

type EncryptionService struct {
	service secret.SecretService
}

func EncryptData(key_encryption, key string) (string, error){
	encription, err := NewEncryptionService(key_encryption)
	if err != nil {
		return "", err
	}
	return encription.Encrypt(key)
}

func DecryptData(key_encryption, key string) (string, error){
	encription, err := NewEncryptionService(key_encryption)
	if err != nil {
		return "", err
	}
	return encription.Decrypt(key)
}

// Returns encription service initialized with a given key
func NewEncryptionService(encodedKey string) (*EncryptionService, error) {
	key, err := secret.EncodedStringToKey(encodedKey)
	if err != nil {
		return nil, err
	}
	s, err := secret.New(&secret.Config{KeyBytes: key})
	if err != nil {
		return nil, err
	}
	return &EncryptionService{
		service: s,
	}, nil
}

// Encrypt a message of in a single string format
func (s *EncryptionService) Encrypt(message string) (string, error) {
	// seal message
	sealed, err := s.service.Seal([]byte(message))
	if err != nil {
		return message, err
	}

	// optionally base64 encode them and store them somewhere (like in a database)
	ciphertext := base64.StdEncoding.EncodeToString(sealed.CiphertextBytes())
	nonce := base64.StdEncoding.EncodeToString(sealed.NonceBytes())

	return nonce + ciphertext, nil
}

// Opens a message from a single string format
func (s *EncryptionService) Decrypt(message string) (string, error) {
	// read in ciphertext and nonce

	if len(message) < 32 {
		return message, errors.New("No nonce")
	}
	nonce, err := base64.StdEncoding.DecodeString(message[:32])
	if err != nil {
		return message, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(message[32:])
	if err != nil {
		return message, err
	}

	// decrypt and open message
	plaintext, err := s.service.Open(&secret.SealedBytes{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	})
	if err != nil {
		return message, err
	}

	return string(plaintext), nil
}

