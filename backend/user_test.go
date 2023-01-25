package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/handler"
)

var router = SetupRouter()


func TestCurrentUser(t *testing.T) {
	for i := 0; i < 1000; i ++ {
		w := httptest.NewRecorder()
		name := "testCurrentUser" + strconv.Itoa(i)
		password := "testCurrentPass" + strconv.Itoa(i)
		input := handler.UserInput{
			Name:     name,
			PassWord: password,
		}
		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
		router.ServeHTTP(w, req)

		w = httptest.NewRecorder()
		input = handler.UserInput{
			Name:     name,
			PassWord: password,
		}
		jsonInput, _ = json.Marshal(input)
		req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
		router.ServeHTTP(w, req)
		jwtToken := w.Body.String()

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/api/user/currentUser", nil)
		req.Header.Add("Authorization", jwtToken)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	for i := 0; i < 1000; i ++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/currentUser", nil)
		router.ServeHTTP(w, req)
		fmt.Println(w.Code)
		assert.NotEqual(t, http.StatusOK, w.Code)
		
	}
}

func TestLogin(t *testing.T) {

	//passwordが入力されていない場合 400
	w := httptest.NewRecorder()
	input := handler.UserInput{
		Name:     "testUser",
		PassWord: "",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// OKな場合 200
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "loginTestUserOK",
		PassWord: "loginTestPassOK",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "loginTestUserOK",
		PassWord: "loginTestPassOK",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, 200, w.Code)

	// 存在しないNameだった場合 400
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "notExistName",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// passwordが間違っている場合 400
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testWrongUser",
		PassWord: "testWrongPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testWrongUser",
		PassWord: "wrongPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

}

func TestSignUp(t *testing.T) {
	//正常な場合 200
	w := httptest.NewRecorder()
	input := handler.UserInput{
		Name:     "testUser",
		PassWord: "testPass",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Passwordがユニークでない場合 200
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testUserUnique1",
		PassWord: "testPassUnique",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testUserUnique2",
		PassWord: "testPassUnique",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Nameがない場合 400
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Passwordがない場合 400
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testUser",
		PassWord: "",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// bodyがない場合 400
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/signUp", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Nameがユニークでない場合 400
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testUserUnique",
		PassWord: "testPassUnique1",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = handler.UserInput{
		Name:     "testUserUnique",
		PassWord: "testPassUnique2",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

}
