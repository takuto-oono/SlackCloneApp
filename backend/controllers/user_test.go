package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
	"backend/controllerUtils"
)

var router = SetupRouter()

type BadResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserId   uint32 `json:"user_id"`
	Username string `json:"username"`
}

func signUpTestFunc(name, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	input := controllerUtils.SignUpAndLoginInput{
		Name:     name,
		Password: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	return w
}

func loginTestFunc(name, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	input := controllerUtils.SignUpAndLoginInput{
		Name:     name,
		Password: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(w, req)
	return w
}

// func TestCurrentUser(t *testing.T) {
// 	errSlice := []int{}
// 	for i := 1; i < 1000; i++ {
// 		fmt.Println(i)
// 		w := httptest.NewRecorder()
// 		name := "testCurrentUser" + strconv.Itoa(i)
// 		password := "testCurrentPass" + strconv.Itoa(i)
// 		input := UserInput{
// 			Name:     name,
// 			PassWord: password,
// 		}
// 		jsonInput, _ := json.Marshal(input)
// 		req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
// 		router.ServeHTTP(w, req)

// 		w = httptest.NewRecorder()
// 		input = UserInput{
// 			Name:     name,
// 			PassWord: password,
// 		}
// 		jsonInput, _ = json.Marshal(input)
// 		req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
// 		router.ServeHTTP(w, req)
// 		jwtToken := w.Body.String()

// 		w = httptest.NewRecorder()
// 		req, _ = http.NewRequest("GET", "/api/user/currentUser", nil)
// 		req.Header.Add("Authorization", jwtToken)
// 		router.ServeHTTP(w, req)
// 		assert.Equal(t, http.StatusOK, w.Code)
// 		fmt.Println(w.Body)
// 		byteArray, _ := ioutil.ReadAll(w.Body)
// 		jsonBody := ([]byte)(byteArray)
// 		user := new(models.User)
// 		if err := json.Unmarshal(jsonBody, user); err != nil {
// 			errSlice = append(errSlice, i)
// 			fmt.Println(err.Error())
// 			continue
// 		}
// 		assert.Equal(t, input.Name, user.Name)
// 		assert.Equal(t, input.PassWord, user.PassWord)
// 	}
// 	fmt.Println(errSlice)

// 	for i := 0; i < 100000; i++ {
// 		w := httptest.NewRecorder()
// 		req, _ := http.NewRequest("GET", "/api/user/currentUser", nil)
// 		router.ServeHTTP(w, req)
// 		fmt.Println(w.Code)
// 		assert.Equal(t, http.StatusBadRequest, w.Code)
// 	}
// }

func TestLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1 正常な場合 200
	// 2 usernameが同一のuserをpasswordで区別できているか 200
	// 3 usernameかpasswordのどちらかがbodyに含まれていない場合 400
	// 4 usernameとpasswordが一致するユーザーが存在しない場合 400

	// 1
	t.Run("1", func(t *testing.T) {
		ids := make([]uint32, 10)
		names := make([]string, 10)

		for i := 0; i < 10; i++ {
			names[i] = "loginControllerTestUser1" + strconv.Itoa(i)
		}

		for i := 0; i < 10; i++ {
			rr := signUpTestFunc(names[i], "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			u := new(models.User)
			json.Unmarshal(jsonBody, u)
			ids[i] = u.ID
			assert.Equal(t, names[i], u.Name)
		}

		for i := 0; i < 10; i++ {
			rr := loginTestFunc(names[i], "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			lr := new(LoginResponse)
			json.Unmarshal(jsonBody, lr)
			assert.NotEmpty(t, lr.Token)
			assert.Equal(t, names[i], lr.Username)
			assert.Equal(t, ids[i], lr.UserId)

		}
	})
	t.Run("2", func(t *testing.T) {
		ids := make([]uint32, 10)
		name := "testLoginUser2"
		passwords := make([]string, 10)
		for i := 0; i < 10; i++ {
			passwords[i] = "TestLoginPass2" + strconv.Itoa(i)
		}
		for i, pass := range passwords {
			rr := signUpTestFunc(name, pass)
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			u := new(models.User)
			json.Unmarshal(jsonBody, u)
			ids[i] = u.ID
			assert.Equal(t, name, u.Name)
		}
		for i, pass := range passwords {
			rr := loginTestFunc(name, pass)
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			lr := new(LoginResponse)
			json.Unmarshal(jsonBody, lr)
			assert.NotEmpty(t, lr.Token)
			assert.Equal(t, name, lr.Username)
			assert.Equal(t, ids[i], lr.UserId)
		}
	})

	t.Run("3", func(t *testing.T) {
		names := make([]string, 10)
		for i := 0; i < 10; i++ {
			name := "testLoginUser3" + strconv.Itoa(i)
			names[i] = name
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}
		for i := 0; i < 10; i++ {
			var rr *httptest.ResponseRecorder
			if i%3 == 0 {
				rr = loginTestFunc("", "")
			} else if i%2 == 1 {
				rr = loginTestFunc(names[i], "")
			} else {
				rr = loginTestFunc("", "pass")
			}
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("4", func(t *testing.T) {
		names := make([]string, 10)
		for i := 0; i < 10; i++ {
			name := "testLoginUser4" + strconv.Itoa(i)
			if i%3 == 0 {
				assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			}
			names[i] = name
		}
		for i, name := range names {
			if i%3 == 0 {
				assert.Equal(t, http.StatusBadRequest, loginTestFunc(name, "wrongPass").Code)
			} else if i%3 == 1 {
				assert.Equal(t, http.StatusBadRequest, loginTestFunc(name, "wrongPass").Code)
			} else {
				assert.Equal(t, http.StatusBadRequest, loginTestFunc(name, "pass").Code)
			}
		}
	})
}

func TestSignUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

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
		for i := 0; i < 12; i++ {
			username := "testSignUpControllerUser3" + strconv.Itoa(i)
			password := "pass"

			var rr *httptest.ResponseRecorder
			if i%3 == 0 {
				rr = signUpTestFunc("", "")
			} else if i%3 == 1 {
				rr = signUpTestFunc(username, "")
			} else {
				rr = signUpTestFunc("", password)
			}
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})

}
