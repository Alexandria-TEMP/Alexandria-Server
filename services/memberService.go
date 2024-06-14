package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type MemberService struct {
	MemberRepository database.ModelRepositoryInterface[*models.Member]
}

func (memberService *MemberService) GetMember(memberID uint) (*models.Member, error) {
	// get member by this id
	member, err := memberService.MemberRepository.GetByID(memberID)
	return member, err
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm, userFields *models.ScientificFieldTagContainer) (*models.LoggedInMemberDTO, error) {
	// check if user with this email already exists
	duplicateMember, err := memberService.MemberRepository.Query(&models.Member{Email: form.Email})
	if err != nil {
		return nil, fmt.Errorf("failed to find all existing member with email %s", form.Email)
	}

	if len(duplicateMember) != 0 {
		return nil, fmt.Errorf("a user with the email %s already exists", form.Email)
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// create member
	// for now no input sanitization for the strings - so first name, last name, email, institution, etc.
	member := &models.Member{
		FirstName:                   form.FirstName,
		LastName:                    form.LastName,
		Email:                       form.Email,
		Password:                    string(passwordHash),
		Institution:                 form.Institution,
		ScientificFieldTagContainer: *userFields,
	}

	// save member to db
	err = memberService.MemberRepository.Create(member)
	if err != nil {
		return nil, err
	}

	// generate tokens
	accessToken, RefreshToken, err := memberService.generateTokenPair(member.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// create logged in member dto
	loggedInMember := &models.LoggedInMemberDTO{
		Member:       member.IntoDTO(),
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
	}

	return loggedInMember, err
}

func (memberService *MemberService) UpdateMember(memberDTO *models.MemberDTO, userFields *models.ScientificFieldTagContainer) error {
	oldMember, err := memberService.MemberRepository.GetByID(memberDTO.ID)
	if err != nil {
		return err
	}

	oldContainer := oldMember.ScientificFieldTagContainer
	oldContainer.ScientificFieldTags = userFields.ScientificFieldTags

	newMember := &models.Member{
		FirstName:                   memberDTO.FirstName,
		LastName:                    memberDTO.LastName,
		Email:                       memberDTO.Email,
		Institution:                 memberDTO.Institution,
		ScientificFieldTagContainer: oldContainer,
	}

	newMember.ID = memberDTO.ID
	_, err = memberService.MemberRepository.Update(newMember)

	return err
}

func (memberService *MemberService) DeleteMember(memberID uint) error {
	err := memberService.MemberRepository.Delete(memberID)
	return err
}

func (memberService *MemberService) GetAllMembers() ([]*models.MemberShortFormDTO, error) {
	members, err := memberService.MemberRepository.Query()

	shortFormDTOs := make([]*models.MemberShortFormDTO, len(members))
	for i, member := range members {
		shortFormDTOs[i] = &models.MemberShortFormDTO{
			ID:        member.ID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
		}
	}

	return shortFormDTOs, err
}

func (memberService *MemberService) LogInMember(form *forms.MemberAuthForm) (*models.LoggedInMemberDTO, error) {
	// get member
	members, err := memberService.MemberRepository.Query(&models.Member{Email: form.Email})
	if err != nil {
		return nil, fmt.Errorf("failed to query members with email %s: %w", form.Email, err)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("there are no members with the email %s", form.Email)
	}

	member := members[0]

	// compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(form.Password)); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// generate tokens
	accessToken, RefreshToken, err := memberService.generateTokenPair(member.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// create logged in member dto
	loggedInMember := &models.LoggedInMemberDTO{
		Member:       member.IntoDTO(),
		AccessToken:  accessToken,
		RefreshToken: RefreshToken,
	}

	return loggedInMember, nil
}

// Credit: https://github.com/war1oc/jwt-auth/blob/master/handler.go
func (memberService *MemberService) RefreshToken(form *forms.TokenRefreshForm) (*models.TokenPairDTO, error) {
	// get secret
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file")
	}

	secret := envFile["SECRET"]

	// get token
	token, err := jwt.Parse(form.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// verify that the typ is refresh
	if token.Header["typ"] != "refresh" {
		return nil, fmt.Errorf("this token is not a refresh token")
	}

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// verify that user exists
		memberID := uint(claims["sub"].(float64))
		if _, err := memberService.MemberRepository.GetByID(memberID); err != nil {
			return nil, fmt.Errorf("the member associated with this token with id=%v doesnt exist: %w", memberID, err)
		}

		// create new access and refresh tokens
		accessToken, refreshToken, err := memberService.generateTokenPair(memberID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token pair: %w", err)
		}

		// return tokens
		return &models.TokenPairDTO{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateTokenPair generates an access token and a refresh token for a member (assumes member is valid and exists)
// The access token is short lived (15 mins), but the refresh token has a long expiration time (3 days).
// When the access token expires, the refresh token can be used to generate a new pair of tokens.
// Credit: https://medium.com/monstar-lab-bangladesh-engineering/jwt-auth-in-go-dde432440924
func (memberService *MemberService) generateTokenPair(memberID uint) (string, string, error) {
	// get secret
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", "", fmt.Errorf("failed to read .env file")
	}

	secret := envFile["SECRET"]

	// CREATE ACCESS TOKEN
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "access"

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = memberID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // 15 min timout

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// CREATE REFRESH TOKEN
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Header["typ"] = "refresh"

	// Set claims
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = memberID
	rtClaims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 day timout

	// Generate encoded token and send it as response.
	rt, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return t, rt, nil
}
