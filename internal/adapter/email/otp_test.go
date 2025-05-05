package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
)

func TestParseOTPTemplate(t *testing.T) {
	tmpl, err := parseOTPTemplate(map[string]any{
		"CompanyName": constant.BankYayaCompanyName,
		"Name":        "Oyen",
		"Purpose":     "login",
		"Code":        "123456",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	var otpHtml = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>This is your OTP for login</title>
</head>
<body style="max-width: 1024px;margin: 0 auto">
<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
  <div style="margin:50px auto;width:70%;padding:20px 0">
    <div style="border-bottom:1px solid #eee">
      <a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">
          PT. Bank Yaya Sumber Uang Tiada Tara (Persero)
      </a>
    </div>
    <p style="font-size:1.1em">Hi, Oyen!</p>
    <p>This is your OTP for login</p>
    <h1>123456</h1>
    <p style="font-size:0.9em;">Regards,<br/>PT. Bank Yaya Sumber Uang Tiada Tara (Persero)</p>
    <hr style="border:none;border-top:1px solid #eee"/>
    <div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
      <p>PT. Bank Yaya Sumber Uang Tiada Tara (Persero)</p>
      <p>Jakarta</p>
      <p>Indonesia</p>
    </div>
  </div>
</div>
</body>
</html>`)

	assert.Equal(t, otpHtml, tmpl)
}
