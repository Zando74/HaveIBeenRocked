package http_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	http_controller "github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/controller/http"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/application/repository/postgresql"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/entity"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/factory"
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/integration/mock"
)

var (
	PasswordList = [][]byte{
		[]byte("p@ssw0rd!"),
		[]byte("Passw0rd#123"),
		[]byte("123456@bcd"),
		[]byte("!Q2w#E4r"),
		[]byte("test1234$"),
		[]byte("9ublic$foo>"),
		[]byte("%admin$p@ss"),
	}
)

func TestMain(m *testing.M) {
	err := mock.SetupTestPostgresqlDB()
	if err != nil {
		log.Fatal(err)
	}
	postgresql.Init()

	m.Run()
}

func TestCheckPasswordController(t *testing.T) {
	router := InitTestServer()

	hashedPasswords := make([]*entity.LeakedHash, len(PasswordList))
	entityFactory := &factory.LeakedHashFactory{}

	for i, password := range PasswordList {
		hashedPasswords[i], _ = entityFactory.Build(password)
	}

	LeakedHashRepositoryImpl.SaveBatch(hashedPasswords)

	for _, password := range PasswordList {

		w := httptest.NewRecorder()

		passwordCheck := http_controller.CheckPasswordValidator{Password: string(password)}
		payload, _ := json.Marshal(passwordCheck)
		req, _ := http.NewRequest("POST", "/api/check", strings.NewReader(string(payload)))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status code 200")

		var response http_controller.CheckPasswordResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "error unmarshaling response")
		assert.True(t, response.Found, "expected password found message")
	}

	w := httptest.NewRecorder()
	unexistingPassword := "unexistingPassword"
	passwordCheck := http_controller.CheckPasswordValidator{Password: unexistingPassword}
	payload, _ := json.Marshal(passwordCheck)
	req, _ := http.NewRequest("POST", "/api/check", strings.NewReader(string(payload)))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "expected status code 200")

	var response http_controller.CheckPasswordResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err, "error unmarshaling response")
	assert.False(t, response.Found, "expected password to not be found")
}
