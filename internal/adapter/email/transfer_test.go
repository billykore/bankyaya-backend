package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
)

func TestParseTransferTemplate(t *testing.T) {
	tmpl, err := parseTransferTemplate(map[string]any{
		"CompanyName":        constant.CompanyName,
		"SourceName":         "Oyen",
		"SourceAccount":      "12345",
		"DestinationName":    "Chiko",
		"DestinationAccount": "54321",
		"DestinationBank":    constant.CompanyName,
		"Amount":             50000,
		"Fee":                0,
		"TransactionRef":     "trf-00001",
		"Note":               "test",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, tmpl)

	otpHtml := []byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>This is your transfer receipt</title>
</head>
<body style="max-width: 1024px;margin: 0 auto">
<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
  <div style="margin:50px auto;width:70%;padding:20px 0">
    <div style="border-bottom:1px solid #eee">
      <a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">
          PT. Bank Yaya Sumber Uang Tiada Tara (Persero), Tbk.
      </a>
    </div>
    <p style="font-size:1.1em">Hi, Oyen!</p>
    <p>This is your transfer receipt</p>
    <p>Source Account: 12345</p>
    <p>To: Chiko</p>
    <p>Destination Account: 54321</p>
    <p>Destination Bank: PT. Bank Yaya Sumber Uang Tiada Tara (Persero), Tbk.</p>
    <p>Amount: 50000</p>
    <p>Fee: 0</p>
    <p>Transaction ID: trf-00001</p>
    <p>Note: test</p>
    <p style="font-size:0.9em;">Regards,<br/>PT. Bank Yaya Sumber Uang Tiada Tara (Persero), Tbk.</p>
    <hr style="border:none;border-top:1px solid #eee"/>
    <div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
      <p>PT. Bank Yaya Sumber Uang Tiada Tara (Persero), Tbk.</p>
      <p>Jakarta</p>
      <p>Indonesia</p>
    </div>
  </div>
</div>
</body>
</html>`)

	assert.Equal(t, otpHtml, tmpl)
}
