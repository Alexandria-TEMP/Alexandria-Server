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
	exampleVersion models.Version
	cwd            string

	memberService        MemberService
	mockTagService       *mocks.MockTagService
	exampleMember        models.Member
	exampleSTag1         *tags.ScientificFieldTag
	exampleSTag2         *tags.ScientificFieldTag
	mockMemberRepository *mocks.MockRepositoryInterface[*models.Member]
)

func TestMain(m *testing.M) {
	exampleVersion = models.Version{
		Model:        gorm.Model{ID: 0},
		Discussions:  nil,
		RenderStatus: models.Pending,
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
