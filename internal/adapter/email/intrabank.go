package email

import (
	"bytes"
	"context"
	"html/template"

	"go.bankyaya.org/app/backend/internal/domain/intrabank"
	"go.bankyaya.org/app/backend/internal/pkg/constant"
	"go.bankyaya.org/app/backend/internal/pkg/email/mailtrap"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
)

var intrabankTmpl = `<!DOCTYPE html>
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
          {{.CompanyName}}
      </a>
    </div>
    <p style="font-size:1.1em">Hi, {{.SourceName}}!</p>
    <p>This is your transfer receipt</p>
    <p>Source Account: {{.SourceAccount}}</p>
    <p>To: {{.DestinationName}}</p>
    <p>Destination Account: {{.DestinationAccount}}</p>
    <p>Destination Bank: {{.DestinationBank}}</p>
    <p>Amount: {{.Amount}}</p>
    <p>Fee: {{.Fee}}</p>
    <p>Transaction ID: {{.TransactionRef}}</p>
    <p>Note: {{.Note}}</p>
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

type IntrabankEmail struct {
	log    *logger.Logger
	client *mailtrap.Client
}

func NewTransferEmail(log *logger.Logger, client *mailtrap.Client) *IntrabankEmail {
	return &IntrabankEmail{
		log:    log,
		client: client,
	}
}

func (e *IntrabankEmail) SendReceipt(_ context.Context, data *intrabank.EmailData) error {
	body, err := parseIntrabankTemplate(map[string]any{
		"CompanyName":        constant.BankYayaCompanyName,
		"SourceName":         data.SourceName,
		"SourceAccount":      data.SourceAccount,
		"DestinationName":    data.DestinationName,
		"DestinationAccount": data.DestinationAccount,
		"DestinationBank":    data.DestinationBank,
		"Amount":             data.Amount,
		"Fee":                data.Fee,
		"TransactionRef":     data.TransactionRef,
		"Note":               data.Note,
	})
	if err != nil {
		e.log.Errorf("SendReceipt error: %v", err)
		return err
	}
	err = e.client.Send(mailtrap.Data{
		Recipient: data.Recipient,
		Subject:   data.Subject,
		Body:      body,
	})
	if err != nil {
		e.log.Errorf("SendReceipt error: %v", err)
		return err
	}
	return nil
}

// buffer to write the email template bytes.
var intrabankTmplBuf = new(bytes.Buffer)

// parseTransferTemplate generates an email template with provided data.
// It returns the generated template as a byte slice or an error if template execution fails.
func parseIntrabankTemplate(data map[string]any) ([]byte, error) {
	defer intrabankTmplBuf.Reset()
	tmpl := template.Must(template.New("intrabank").Parse(intrabankTmpl))
	err := tmpl.Execute(intrabankTmplBuf, data)
	if err != nil {
		return nil, err
	}
	return intrabankTmplBuf.Bytes(), nil
}
