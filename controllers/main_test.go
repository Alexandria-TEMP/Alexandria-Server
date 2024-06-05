package controllers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/services/interfaces"
)

var (
	cwd              string
	router           *gin.Engine
	responseRecorder *httptest.ResponseRecorder

	branchController  BranchController
	mockBranchService interfaces.BranchService
)

// TestMain is a keyword function, this is run by the testing package before other tests
func TestMain(m *testing.M) {
	// Setup test router, to test controller endpoints through http
	router = SetUpRouter()

	cwd, _ = os.Getwd()

	os.Exit(m.Run())
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router = gin.Default()

	return router
}
