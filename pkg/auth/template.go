package auth

import (
	"net/url"
	"strings"

	"github.com/sev-2/raiden"
)

type MailClient interface {
	Mail(string, string, string, string, map[string]interface{}) error
}

// TemplateMailer will send mail and use templates from the site for easy mail styling
type TemplateMailer struct {
	Config *raiden.Config
	Mailer MailClient
}

func encodeRedirectURL(referrerURL string) string {
	if len(referrerURL) > 0 {
		if strings.ContainsAny(referrerURL, "&=#") {
			referrerURL = url.QueryEscape(referrerURL)
		}
	}
	return referrerURL
}

const defaultRecoveryMail = `<h2>Reset password</h2>

<p>Follow this link to reset the password for your user:</p>
<p><a href="{{ .ConfirmationURL }}">Reset password</a></p>
<p>Alternatively, enter the code: {{ .Token }}</p>`

// RecoveryMail sends a password recovery mail
func (m *TemplateMailer) RecoveryMail(email string, token string, otp, referrerURL string, externalURL *url.URL) error {
	path, err := getPath("/verify", &EmailParams{
		Token:      token,
		Type:       "recovery",
		RedirectTo: referrerURL,
	})
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"ConfirmationURL": externalURL.ResolveReference(path).String(),
		"Email":           email,
		"Token":           otp,
		"TokenHash":       token,
		"RedirectTo":      referrerURL,
	}

	return m.Mailer.Mail(
		email,
		"Reset Your Password",
		"",
		defaultRecoveryMail,
		data,
	)
}

// Send mail function
func (m *TemplateMailer) Send(email, subject, body string, data map[string]interface{}) error {
	return m.Mailer.Mail(
		email,
		subject,
		"",
		body,
		data,
	)
}
