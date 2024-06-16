package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"go.uber.org/mock/gomock"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	mockMemberService *mocks.MockMemberService
	middleware        Middleware
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.GET("/auth", middleware.CheckAuth, func(c *gin.Context) { c.Status(http.StatusOK) })

	cwd, _ = os.Getwd()

	code := m.Run()
	os.Exit(code)
}

func beforeEach(t *testing.T) {
	t.Helper()

	responseRecorder = httptest.NewRecorder()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMemberService = mocks.NewMockMemberService(mockCtrl)
	middleware = Middleware{
		MemberService: mockMemberService,
		Secret:        "secret",
	}
}

func TestCheckAuthSuccess(t *testing.T) {
	beforeEach(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(nil, nil)

	validToken, err := generateAccessToken(t, "access")
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode)
}

func TestCheckAuthMissingAuth(t *testing.T) {
	beforeEach(t)

	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestCheckAuthBadHeaderFormat(t *testing.T) {
	beforeEach(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(nil, nil)

	validToken, err := generateAccessToken(t, "access")
	assert.Nil(t, err)

	// Bear instead of Bearer
	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bear %s", validToken))
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)

	// Extra space
	req, _ = http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bear %s a", validToken))
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestCheckAuthWontParse(t *testing.T) {
	beforeEach(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(nil, nil)

	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOidIUzI1NiJ9.eyJleHAiOjE3MTg1MzY5MDQsInN1YiI6MTF9.V-r-uUYQnSPL1A6k2YmorOI1eRz5Ah2QRoLGuGyDj30"))
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestCheckAuthNotAccessTyp(t *testing.T) {
	beforeEach(t)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(nil, nil)

	token, err := generateAccessToken(t, "refresh")
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func TestCheckAuthExpired(t *testing.T) {
	beforeEach(t)

	token, err := generateAccessToken(t, "access")
	assert.Nil(t, err)
	time.Sleep(time.Second * 7)

	mockMemberService.EXPECT().GetMember(uint(1)).Return(nil, nil)

	req, _ := http.NewRequest("GET", "/auth", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	router.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode)
}

func generateAccessToken(t *testing.T, typ string) (string, error) {
	t.Helper()

	// CREATE ACCESS TOKEN
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = typ

	// Set claims
	claims, _ := token.Claims.(jwt.MapClaims)
	claims["sub"] = uint(1)
	claims["exp"] = time.Now().Add(time.Second * 5).Unix() // 5 sec timeout

	stringToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return stringToken, err
}
