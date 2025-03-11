package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/sev-2/raiden/pkg/logger"
)

var RecoverLogger = logger.HcLog().Named("recover")

type RecoverResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Otp         string `json:"otp"`
}

type AuthExtension struct {
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

func (auth *AuthExtension) Recover(ctx context.Context, email string, referrerURL string) (*RecoverResponse, error) {
	otp, err := generateOtp(6)
	if err != nil {
		RecoverLogger.Error("error generate otp", "err", err)
		return nil, err
	}

	token := generateTokenHash(email, otp)
	return &RecoverResponse{
		AccessToken: token,
		TokenType:   "recovery",
		Otp:         otp,
	}, nil
}
