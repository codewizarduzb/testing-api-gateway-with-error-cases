package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing-api-gateway/api_test/handlers"
	"testing-api-gateway/api_test/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	gin.SetMode(gin.TestMode)
	require.NoError(t, SetupMinimumInstance(""))
	buffer, err := OpenFile("user.json")

	require.NoError(t, err)
	// User Create
	req := NewRequest(http.MethodPost, "/user/create", buffer)
	res := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/user/create", handlers.CreateUser)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)

	var user storage.User
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &user))

	require.Equal(t, "xqoshmaqboyev@gmail.com", user.Email)
	require.Equal(t, int64(21), user.Age)
	require.Equal(t, "Xumoyunmirzo", user.FirstName)
	require.Equal(t, "codewizarduzbekistan", user.Username)
	require.Equal(t, "12345", user.Password)
	require.NotNil(t, user.Id)

	// GetUser
	getReq := NewRequest(http.MethodGet, "/users/get", buffer)
	q := getReq.URL.Query()
	q.Add("id", user.Id)
	getReq.URL.RawQuery = q.Encode()
	getRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users/get", handlers.GetUser)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)
	var getResp storage.User
	userBytes, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(userBytes, &getResp))
	assert.Equal(t, user.Id, getResp.Id)
	assert.Equal(t, user.FirstName, getResp.FirstName)
	assert.Equal(t, user.Username, getResp.Username)
	assert.Equal(t, user.Age, getResp.Age)
	assert.Equal(t, user.Email, getResp.Email)

	// Users List
	listReq := NewRequest(http.MethodGet, "/users", buffer)
	listRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users", handlers.ListUsers)
	r.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)
	userBytes, err = io.ReadAll(listRes.Body)
	assert.NoError(t, err)
	assert.NotNil(t, userBytes)

	// User Register
	regReq := NewRequest(http.MethodPost, "/user/register", buffer)
	regRes := httptest.NewRecorder()
	r.POST("/user/register", handlers.RegisterUser)
	r.ServeHTTP(regRes, regReq)
	assert.Equal(t, 500, regRes.Code)
	var resp storage.Message
	userBytes, err = io.ReadAll(regRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(userBytes, &resp))
	require.NotNil(t, resp.Message)
	require.Equal(t, "", resp.Message)

	// User Verify with correct code
	verifyURL := "/user/verify/12345"
	verifyReq := NewRequest(http.MethodGet, verifyURL, buffer)
	verifyResp := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/user/verify/:code", handlers.Verify)
	r.ServeHTTP(verifyResp, verifyReq)

	assert.Equal(t, http.StatusOK, verifyResp.Code)
	var respCorrect storage.Message
	bodyBytesCorrect, err := io.ReadAll(verifyResp.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytesCorrect, &respCorrect))
	require.Equal(t, "Success", respCorrect.Message)


	// ***************************** TESTS FOR ERROR CASES ************************
	
	// Test for user failure
	failingBuffer := []byte(`{"invalid": "json"}`)
	failCreateReq := NewRequest(http.MethodPost, "/user/create", failingBuffer)
	failCreateRes := httptest.NewRecorder()
	r.ServeHTTP(failCreateRes, failCreateReq)
	assert.Equal(t, 404, failCreateRes.Code)

	// Test for failure in getting user information
	failGetReq := NewRequest(http.MethodGet, "/users/get", buffer)
	q = failGetReq.URL.Query()
	q.Add("id", "nonexistent_id")
	failGetReq.URL.RawQuery = q.Encode()
	failGetRes := httptest.NewRecorder()
	r.ServeHTTP(failGetRes, failGetReq)
	assert.Equal(t, http.StatusNotFound, failGetRes.Code)

	// Test for failure in listing users
	failListReq := NewRequest(http.MethodGet, "/users", nil)
	failListRes := httptest.NewRecorder()
	r.ServeHTTP(failListRes, failListReq)
	assert.Equal(t, 404, failListRes.Code)

	// Test for failure in user registration
	failRegReq := NewRequest(http.MethodPost, "/user/register", nil)
	failRegRes := httptest.NewRecorder()
	r.ServeHTTP(failRegRes, failRegReq)
	assert.Equal(t, 404, failRegRes.Code)

	// Test for failure in user verification
	failVerifyReq := NewRequest(http.MethodGet, "/user/verify/invalid_code", nil)
	failVerifyRes := httptest.NewRecorder()
	r.ServeHTTP(failVerifyRes, failVerifyReq)
	assert.Equal(t, http.StatusBadRequest, failVerifyRes.Code)
}
