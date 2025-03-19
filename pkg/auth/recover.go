package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"strconv"

	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/pkg/logger"
)

func NewLibrary(config *raiden.Config) any {
	return &AuthExtension{
		config: config,
	}
}

type AuthExtension struct {
	raiden.BaseLibrary
	config *raiden.Config
}

var RecoverLogger = logger.HcLog().Named("recover")

type RecoverResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Otp         string `json:"otp"`
}

// GenerateOtp generates a random n digit otp
func generateOtp(digits int) (string, error) {
	upper := math.Pow10(digits)
	val, err := rand.Int(rand.Reader, big.NewInt(int64(upper)))
	if err != nil {
		return "", fmt.Errorf("error generate otp : %s", err)
	}
	// adds a variable zero-padding to the left to ensure otp is uniformly random
	expr := "%0" + strconv.Itoa(digits) + "v"
	otp := fmt.Sprintf(expr, val.String())
	return otp, nil
}

func generateTokenHash(emailOrPhone, otp string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(emailOrPhone+otp)))
}

func (auth *AuthExtension) Recover(email string, referrerURL string) error {
	otp, err := generateOtp(6)
	if err != nil {
		RecoverLogger.Error("error generate otp", "err", err)
		return err
	}

	user, err := GetUserByEmail(auth.config, email)
	if err != nil {
		RecoverLogger.Error("error get user", "err", err)
		return err
	}

	token := generateTokenHash(user.Email, otp)

	// send recovery mail
	mailer := NewMailer(auth.config)

	url, err := url.Parse(referrerURL)
	if err != nil {
		RecoverLogger.Error("error parse referrer url", "err", err)
		return err
	}

	err = mailer.RecoveryMail(user.Email, token, otp, referrerURL, url)
	if err != nil {
		RecoverLogger.Error("error send recovery mail", "err", err)
		return err
	}

	err = UpdateUserRecoveryToken(auth.config, email, token)
	if err != nil {
		RecoverLogger.Error("error update user recovery token", "err", err)
		return err
	}

	return nil
}
