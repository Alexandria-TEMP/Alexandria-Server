package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

var (
	versionService VersionService
	c              *gin.Context
	mockFilesystem *mocks.MockFilesystem
	cwd            string

	memberService         MemberService
	exampleMember         models.Member
	exampleSTag1          *tags.ScientificFieldTag
	exampleSTag2          *tags.ScientificFieldTag
	mockMemberRepository  *mocks.MockRepositoryInterface[*models.Member]
	mockVersionRepository *mocks.MockRepositoryInterface[*models.Version]

	pendingVersion models.Version
	failureVersion models.Version
	successVersion models.Version
)

func TestMain(m *testing.M) {
	pendingVersion = models.Version{
		Model:        gorm.Model{ID: 0},
		Discussions:  nil,
		RenderStatus: models.RenderPending,
	}

	failureVersion = models.Version{
		Model:        gorm.Model{ID: 1},
		Discussions:  nil,
		RenderStatus: models.RenderFailure,
	}

	successVersion = models.Version{
		Model:        gorm.Model{ID: 2},
		Discussions:  nil,
		RenderStatus: models.RenderSuccess,
	}
	tag1 := tags.ScientificFieldTag{
		ScientificField: "Mathematics",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag1 = &tag1
	tag2 := tags.ScientificFieldTag{
		ScientificField: "",
		Subtags:         []*tags.ScientificFieldTag{},
		ParentID:        nil,
	}
	exampleSTag2 = &tag2

	exampleMember = models.Member{
		FirstName:           "John",
		LastName:            "Smith",
		Email:               "john.smith@gmail.com",
		Password:            "password",
		Institution:         "TU Delft",
		ScientificFieldTags: []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2},
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
