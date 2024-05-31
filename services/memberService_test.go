package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockTagService = mocks.NewMockTagService(mockCtrl)
	mockMemberRepository := mocks.NewMockRepositoryInterface[*models.Member](mockCtrl)

	//need to mock repository here, how?
	memberService = MemberService{
		MemberRepository: mockMemberRepository,
	}
}


func testGetMember(t *testing.T, id uint) {
	// mock member repository here to return member when get by id called

	// call service method
	member, err := memberService.GetMember(id)
	// assert member was returned correctly
	assert.Equal(t, exampleMember, member)
	// assert there was no error
	assert.Nil(t, err)

}
