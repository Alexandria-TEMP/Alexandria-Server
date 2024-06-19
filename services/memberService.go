package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccessTokenDuration  = 15 // mins
	RefreshTokenDuration = 72 // hours
)

type MemberService struct {
	Secret           string
	MemberRepository database.ModelRepositoryInterface[*models.Member]
}

func (memberService *MemberService) GetMember(memberID uint) (*models.Member, error) {
	// get member by this id
	member, err := memberService.MemberRepository.GetByID(memberID)
	return member, err
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm, userFields *models.ScientificFieldTagContainer) (string, int64, string, int64, *models.Member, error) {
	// check if user with this email already exists
	duplicateMember, err := memberService.MemberRepository.Query(&models.Member{Email: form.Email})
	if err != nil {
		return "", 0, "", 0, nil, fmt.Errorf("failed to find all existing member with email %s", form.Email)
	}

	if len(duplicateMember) != 0 {
		return "", 0, "", 0, nil, fmt.Errorf("a user with the email %s already exists", form.Email)
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, "", 0, nil, fmt.Errorf("failed to hash password: %w", err)
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
		return "", 0, "", 0, nil, err
	}

	// generate tokens
	accessToken, aExp, refreshToken, rExp, err := memberService.generateTokenPair(member.ID)
	if err != nil {
		return "", 0, "", 0, nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return accessToken, aExp, refreshToken, rExp, member, nil
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

func (memberService *MemberService) LogInMember(form *forms.MemberAuthForm) (*models.Member, string, int64, string, int64, error) {
	// get member
	members, err := memberService.MemberRepository.Query(&models.Member{Email: form.Email})
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("failed to query members with email %s: %w", form.Email, err)
	}

	if len(members) == 0 {
		return nil, "", 0, "", 0, fmt.Errorf("no members with email %s found", form.Email)
	}

	member := members[0]

	// compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(form.Password)); err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("invalid password")
	}

	// generate tokens
	accessToken, aExp, RefreshToken, rExp, err := memberService.generateTokenPair(member.ID)
	if err != nil {
		return nil, "", 0, "", 0, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return member, accessToken, aExp, RefreshToken, rExp, nil
}

// Credit: https://github.com/war1oc/jwt-auth/blob/master/handler.go
func (memberService *MemberService) RefreshToken(form *forms.TokenRefreshForm) (string, int64, string, int64, error) {
	// get token
	token, err := jwt.Parse(form.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(memberService.Secret), nil
	})

	if err != nil {
		return "", 0, "", 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// verify that the typ is refresh
	if token.Header["typ"] != "refresh" {
		return "", 0, "", 0, fmt.Errorf("this token is not a refresh token")
	}

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// verify that user exists
		memberID := uint(claims["sub"].(float64))
		if _, err := memberService.MemberRepository.GetByID(memberID); err != nil {
			return "", 0, "", 0, fmt.Errorf("the member associated with this token with id=%v doesnt exist: %w", memberID, err)
		}

		// create new access and refresh tokens
		accessToken, aExp, refreshToken, rExp, err := memberService.generateTokenPair(memberID)
		if err != nil {
			return "", 0, "", 0, fmt.Errorf("failed to generate token pair: %w", err)
		}

		return accessToken, aExp, refreshToken, rExp, nil
	}

	return "", 0, "", 0, fmt.Errorf("invalid token")
}

// generateTokenPair generates an access token and a refresh token for a member (assumes member is valid and exists)
// The access token is short lived (15 mins), but the refresh token has a long expiration time (3 days).
// When the access token expires, the refresh token can be used to generate a new pair of tokens.
// Credit: https://medium.com/monstar-lab-bangladesh-engineering/jwt-auth-in-go-dde432440924
func (memberService *MemberService) generateTokenPair(memberID uint) (at string, aexp int64, rt string, rexp int64, err error) {
	// CREATE ACCESS TOKEN
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["typ"] = "access"

	// Set claims
	claims, _ := token.Claims.(jwt.MapClaims)
	claims["sub"] = memberID
	aexp = time.Now().Add(time.Minute * time.Duration(AccessTokenDuration)).Unix() // 15 min timout
	claims["exp"] = aexp

	at, err = token.SignedString([]byte(memberService.Secret))
	if err != nil {
		return "", 0, "", 0, err
	}

	// CREATE REFRESH TOKEN
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Header["typ"] = "refresh"

	// Set claims
	rtClaims, _ := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = memberID
	rexp = time.Now().Add(time.Hour * time.Duration(RefreshTokenDuration)).Unix() // 3 day timout
	rtClaims["exp"] = rexp

	rt, err = refreshToken.SignedString([]byte(memberService.Secret))
	if err != nil {
		return "", 0, "", 0, err
	}

	return at, aexp, rt, rexp, nil
}
