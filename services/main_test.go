package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
)

var (
	mockFilesystem *mocks.MockFilesystem

	c   *gin.Context
	cwd string
)

func TestMain(m *testing.M) {
	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
