package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

const (
	ScopeAuth = "authentication"
)

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserId: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	emptyByte := make([]byte, 32)

	_, err := rand.Read(emptyByte)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyByte)

	hash := sha256.Sum256([]byte(token.Plaintext))

	token.Hash = hash[:]

	return token, nil
}
