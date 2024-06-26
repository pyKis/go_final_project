package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pyKis/go_final_project/configs/models"

	"github.com/golang-jwt/jwt"
)

func SignInPOST(w http.ResponseWriter, r *http.Request) {
	var signData models.Sign
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &signData); err != nil {
		responseWithError(w, "JSON encoding error", err)
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")

	if signData.Password == envPassword {
		jwtInstance := jwt.New(jwt.SigningMethodHS256)
		token, err := jwtInstance.SignedString([]byte(envPassword))
log.Println(token)
		taskIdData, err := json.Marshal(models.AuthToken{Token: token})
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(taskIdData)

		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}
	} else {
		errorResponse := models.ErrorResponse{Error: "wrong password"}
		errorData, _ := json.Marshal(errorResponse)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write(errorData)

		if err != nil {
			http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusUnauthorized)
		}
	}
}