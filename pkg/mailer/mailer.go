package mailer

import (
	"fmt"
	"net/url"

	"github.com/sev-2/raiden"
	"github.com/sev-2/raiden/pkg/logger"
	"github.com/supabase/mailme"
	"gopkg.in/gomail.v2"
)

var MailLogger = logger.HcLog().Named("mailer")

type Mailer interface {
	RecoveryMail(email, token, otp, referrerURL string, externalURL *url.URL) error
	Send(email, subject, body string, data map[string]interface{}) error
}

type Config struct {
	SmtpHost       string
	SmtpPort       int
	SmtpUser       string
	SmtpPass       string
	SmtpAdminEmail string
	SmtpSenderName string
}

type EmailParams struct {
	Token      string
	Type       string
	RedirectTo string
}

func NewMailer(config *raiden.Config) Mailer {

	mail := gomail.NewMessage()

	var adminEmail = config.GetString("SMTP_ADMIN_EMAIL")
	var senderName = config.GetString("SMTP_SENDER_NAME")
	var smtpHost = config.GetString("SMTP_HOST")
	var smtpPort = config.GetInt("SMTP_PORT")
	var smtpUser = config.GetString("SMTP_USER")
	var smtpPass = config.GetString("SMTP_PASS")

	from := mail.FormatAddress(adminEmail, senderName)

	var mailClient = &mailme.Mailer{
		Host: smtpHost,
		Port: smtpPort,
		User: smtpUser,
		Pass: smtpPass,
		From: from,
	}

	return &TemplateMailer{
		Config: config,
		Mailer: mailClient,
	}
}

func getPath(filepath string, params *EmailParams) (*url.URL, error) {
	path := &url.URL{}
	if filepath != "" {
		if p, err := url.Parse(filepath); err != nil {
			return nil, err
		} else {
			path = p
		}
	}
	if params != nil {
		path.RawQuery = fmt.Sprintf("token=%s&type=%s&redirect_to=%s", url.QueryEscape(params.Token), url.QueryEscape(params.Type), encodeRedirectURL(params.RedirectTo))
	}
	return path, nil
}
