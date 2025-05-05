package email

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"go.bankyaya.org/app/backend/internal/domain/otp"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
	"go.bankyaya.org/app/backend/internal/pkg/email/mailtrap"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
)

var otpTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>This is your OTP for {{.Purpose}}</title>
</head>
<body style="max-width: 1024px;margin: 0 auto">
<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
  <div style="margin:50px auto;width:70%;padding:20px 0">
    <div style="border-bottom:1px solid #eee">
      <a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">
          {{.CompanyName}}
      </a>
    </div>
    <p style="font-size:1.1em">Hi, {{.Name}}!</p>
    <p>This is your OTP for {{.Purpose}}</p>
    <h1>{{.Code}}</h1>
    <p style="font-size:0.9em;">Regards,<br/>{{.CompanyName}}</p>
    <hr style="border:none;border-top:1px solid #eee"/>
    <div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
      <p>{{.CompanyName}}</p>
      <p>Jakarta</p>
      <p>Indonesia</p>
    </div>
  </div>
</div>
</body>
</html>`

type OTPEmail struct {
	log    *logger.Logger
	client *mailtrap.Client
}

func NewOTPEmail(log *logger.Logger, client *mailtrap.Client) *OTPEmail {
	return &OTPEmail{
		log:    log,
		client: client,
	}
}

func (e *OTPEmail) Send(_ context.Context, otpData *otp.OTP) error {
	if otpData.Channel != otp.ChannelEmail {
		return errors.New("invalid channel")
	}
	body, err := parseOTPTemplate(map[string]any{
		"CompanyName": constant.BankYayaCompanyName,
		"Name":        otpData.User.Name,
		"Purpose":     otpData.Purpose,
		"Code":        otpData.Code,
	})
	if err != nil {
		e.log.Errorf("Send error: %v", err)
		return err
	}
	err = e.client.Send(mailtrap.Data{
		Recipient: otpData.User.Email,
		Subject:   otpData.Purpose.Message(),
		Body:      body,
	})
	if err != nil {
		e.log.Errorf("Send error: %v", err)
		return err
	}
	return nil
}

// buffer to write the email template bytes.
var otpTmplBuf = new(bytes.Buffer)

// parseOTPTemplate generates an email template by populating placeholders with provided data.
// Returns the generated template as a byte slice or an error if template execution fails.
func parseOTPTemplate(data map[string]any) ([]byte, error) {
	defer otpTmplBuf.Reset()
	tmpl := template.Must(template.New("otp").Parse(otpTmpl))
	err := tmpl.Execute(otpTmplBuf, data)
	if err != nil {
		return nil, err
	}
	return otpTmplBuf.Bytes(), nil
}
