package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type JWTPayload struct {
	UserID uuid.UUID
	Exp    int64
}

func GenerateJWT(userId uuid.UUID, secret string) (string, error) {
	payload := JWTPayload{
		UserID: userId,
		Exp:    time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	header := `{"alg":"HS256","typ":"JWT"}`
	data := header + "." + base64.URLEncoding.EncodeToString(payloadBytes)

	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write([]byte(data))
	if err != nil {
		return "", err
	}

	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return header + "." + base64.URLEncoding.EncodeToString(payloadBytes) + "." + signature, nil
}

func verifyJWT(token, secret string) (uuid.UUID, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return uuid.Nil, errors.New("invalid token")
	}

	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return uuid.Nil, err
	}

	var payloadObj JWTPayload
	err = json.Unmarshal(payload, &payloadObj)
	if err != nil {
		return uuid.Nil, err
	}

	data := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write([]byte(data))
	if err != nil {
		return uuid.Nil, err
	}

	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if expectedSignature != parts[2] {
		return uuid.Nil, errors.New("invalid signature")
	}

	return payloadObj.UserID, nil
}
