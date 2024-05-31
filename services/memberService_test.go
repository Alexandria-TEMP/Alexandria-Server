package services

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/mocks"
)

func beforeEachMember(t *testing.T) {
	t.Helper()

	mockCtrl := gomock.NewController(t)

	defer mockCtrl.Finish()

	mockTagService = mocks.NewMockTagService(mockCtrl)

	//need to mock repository here, how?
	memberService = MemberService{}
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
