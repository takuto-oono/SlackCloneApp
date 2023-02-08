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

func signUpTestFunc(name, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	input := models.User{
		Name:     name,
		PassWord: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	return w
}

func loginTestFunc(name, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	input := UserInput{
		Name:     name,
		PassWord: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	return w
}

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
	// w := httptest.NewRecorder()
	// input := UserInput{
	// 	Name:     "testUser",
	// 	PassWord: "",
	// }
	// jsonInput, _ := json.Marshal(input)
	// req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusBadRequest, w.Code)

	noPassWordTestLoginNames := []string{}
	for i := 0; i < 1000; i++ {
		noPassWordTestLoginNames = append(noPassWordTestLoginNames, "noPassWordTestLoginNames"+strconv.Itoa(i))
	}
	for _, name := range noPassWordTestLoginNames {
		rr := signUpTestFunc(name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
	}

	for _, name := range noPassWordTestLoginNames {
		rr := loginTestFunc(name, "")
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	}
	// OKな場合 200
	// w = httptest.NewRecorder()
	// input = UserInput{
	// 	Name:     "loginTestUserOK",
	// 	PassWord: "loginTestPassOK",
	// }
	// jsonInput, _ = json.Marshal(input)
	// req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// w = httptest.NewRecorder()
	// input = UserInput{
	// 	Name:     "loginTestUserOK",
	// 	PassWord: "loginTestPassOK",
	// }
	// jsonInput, _ = json.Marshal(input)
	// req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// fmt.Println(w)
	// assert.Equal(t, http.StatusOK, w.Code)

	// // 存在しないNameだった場合 400
	// w = httptest.NewRecorder()
	// input = UserInput{
	// 	Name:     "notExistName",
	// 	PassWord: "testPass",
	// }
	// jsonInput, _ = json.Marshal(input)
	// req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusBadRequest, w.Code)

	// // passwordが間違っている場合 400
	// w = httptest.NewRecorder()
	// input = UserInput{
	// 	Name:     "testWrongUser",
	// 	PassWord: "testWrongPass",
	// }
	// jsonInput, _ = json.Marshal(input)
	// req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// w = httptest.NewRecorder()
	// input = UserInput{
	// 	Name:     "testWrongUser",
	// 	PassWord: "wrongPass",
	// }
	// jsonInput, _ = json.Marshal(input)
	// req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	// router.ServeHTTP(w, req)
	// assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestSignUp(t *testing.T) {
	// 1 普通の場合 200
	// 2 usernameがuniqueでない場合 200
	// 3 usernameかpasswordがbodyに含まれていない場合 400
	// 4 usernameとpasswordが同一のuserが既に存在している場合 400

	// 1
	t.Run("1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			username := "testSignUpControllerUser1" + strconv.Itoa(i)
			password := "pass"
			assert.Equal(t, http.StatusOK, signUpTestFunc(username, password).Code)
		}
	})

	// 2
	t.Run("2", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			username := "testSignUpControllerUser2"
			password := "password" + strconv.Itoa(i)
			assert.Equal(t, http.StatusOK, signUpTestFunc(username, password).Code)
		}
	})

	// 3
	t.Run("3", func(t *testing.T) {
		for i := 0; i < 12; i ++ {
			username := "testSignUpControllerUser2" + strconv.Itoa(i)
			password := "pass"
			
			var rr *httptest.ResponseRecorder
			if i % 3 == 0 {
				rr = signUpTestFunc("", "")
			} else if i % 3 == 1 {
				rr = signUpTestFunc(username, "")
			} else {
				rr = signUpTestFunc("", password)
			}
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})

}
