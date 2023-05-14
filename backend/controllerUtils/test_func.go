package controllerUtils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"backend/controllerUtils"
	"backend/models"
)

var testRouter = SetupRouter()

func Req(method, endpoint, jwtToken string, bodyInfo interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(bodyInfo)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	if jwtToken != "" {
		req.Header.Set("Authorization", jwtToken)
	}
	testRouter.ServeHTTP(rr, req)
	return rr
}

func Res(rr *httptest.ResponseRecorder, body interface{}) (int, interface{}) {
	byteArray, _ := io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), &body)
	return rr.Code, body
}

func signUpTestFunc(userName, pass string) (*httptest.ResponseRecorder, models.User) {
	rr := Req(http.MethodPost, "/api/user/signUp", "", controllerUtils.SignUpAndLoginInput{
		Name: userName,
		Password: pass,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var u models.User
	json.Unmarshal(([]byte)(byteArray), &u)
	return rr, u
}

func loginTestFunc(name, password string) (*httptest.ResponseRecorder, LoginResponse) {
	rr := Req(http.MethodPost, "/api/user/login", "", controllerUtils.SignUpAndLoginInput{
		Name: name,
		Password: password,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var lr LoginResponse
	json.Unmarshal(([]byte)(byteArray), &lr)
	return rr, lr
}

func createChannelTestFunc(name, description string, isPrivate *bool, jwtToken string, workspaceID int) (*httptest.ResponseRecorder, models.Channel) {

}

func addUserInChannelTestFunc(channelID int, userID uint32, jwtToken string) *httptest.ResponseRecorder, 
