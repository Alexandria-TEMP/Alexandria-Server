package services

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

var (
	c   *gin.Context
	cwd string

	memberService        MemberService
	exampleMember        models.Member
	exampleSTag1         *tags.ScientificFieldTag
	exampleSTag2         *tags.ScientificFieldTag
	mockMemberRepository *mocks.MockRepositoryInterface[*models.Member]
)

func TestMain(m *testing.M) {
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
	scientificFieldTagArray := []*tags.ScientificFieldTag{exampleSTag1, exampleSTag2}
	exampleMember = models.Member{
		FirstName:   "John",
		LastName:    "Smith",
		Email:       "john.smith@gmail.com",
		Password:    "password",
		Institution: "TU Delft",
		ScientificFieldTagContainer: &tags.ScientificFieldTagContainer{
			ScientificFieldTags: scientificFieldTagArray,
		},
	}

	cwd, _ = os.Getwd()

	c, _ = gin.CreateTestContext(httptest.NewRecorder())

	os.Exit(m.Run())
}
