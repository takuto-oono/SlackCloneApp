package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
)

var router = SetupRouter()

func TestCurrentUser(t *testing.T) {
	errSlice := []int{}
	for i := 1; i < 1000; i++ {
		fmt.Println(i)
		w := httptest.NewRecorder()
		name := "testCurrentUser" + strconv.Itoa(i)
		password := "testCurrentPass" + strconv.Itoa(i)
		input := UserInput{
			Name:     name,
			PassWord: password,
		}
		jsonInput, _ := json.Marshal(input)
		req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
		router.ServeHTTP(w, req)

		w = httptest.NewRecorder()
		input = UserInput{
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
		fmt.Println(w.Body)
		byteArray, _ := ioutil.ReadAll(w.Body)
		jsonBody := ([]byte)(byteArray)
		user := new(models.User)
		if err := json.Unmarshal(jsonBody, user); err != nil {
			errSlice = append(errSlice, i)
			fmt.Println(err.Error())
			continue
		}
		assert.Equal(t, input.Name, user.Name)
		assert.Equal(t, input.PassWord, user.PassWord)
	}
	fmt.Println(errSlice)

	for i := 0; i < 100000; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user/currentUser", nil)
		router.ServeHTTP(w, req)
		fmt.Println(w.Code)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestLogin(t *testing.T) {

	//passwordが入力されていない場合 400
	w := httptest.NewRecorder()
	input := UserInput{
		Name:     "testUser",
		PassWord: "",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// OKな場合 200
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "loginTestUserOK",
		PassWord: "loginTestPassOK",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "loginTestUserOK",
		PassWord: "loginTestPassOK",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	fmt.Println(w)
	assert.Equal(t, http.StatusOK, w.Code)

	// 存在しないNameだった場合 400
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "notExistName",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// passwordが間違っている場合 400
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "testWrongUser",
		PassWord: "testWrongPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "testWrongUser",
		PassWord: "wrongPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func signUpTestFunc(name, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	input := UserInput{
		Name:     name,
		PassWord: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	return w
}

func TestSignUp2(t *testing.T) {
	w := signUpTestFunc("testFuncUser", "pass")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignUp(t *testing.T) {
	//正常な場合 200
	w := httptest.NewRecorder()
	input := UserInput{
		Name:     "testUser",
		PassWord: "testPass",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Passwordがユニークでない場合 200
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "testUserUnique1",
		PassWord: "testPassUnique",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	w = httptest.NewRecorder()
	input = UserInput{
		Name:     "testUserUnique2",
		PassWord: "testPassUnique",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Nameがない場合 400
	for i := 0; i < 1000; i ++ {
		rr := signUpTestFunc("", "pass" + strconv.Itoa(i))
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}

	// Passwordがない場合 400
	noPassWordNames := []string{}
	for i := 0; i < 1000; i ++ {
		noPassWordNames = append(noPassWordNames, "noPassWordTestName" + strconv.Itoa(i))
	}
	for _, name := range noPassWordNames {
		rr := signUpTestFunc(name, "")
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
	
	// bodyがない場合 400
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/signUp", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Nameがユニークでない場合 400
	notUniqueNames := []string{}
	for i := 0; i < 1000; i ++ {
		notUniqueNames = append(notUniqueNames, "testUserNameNotUnique" + strconv.Itoa(i))
	}
	for _, name := range notUniqueNames {
		rr := signUpTestFunc(name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr = signUpTestFunc(name, "pass")
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
	
}
