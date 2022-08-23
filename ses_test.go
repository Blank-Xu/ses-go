package ses

import (
	"testing"
)

const (
	endpoint        = "endpoint"
	from            = "from"
	accessKeyID     = "accessKeyID"
	secretAccessKey = "secretAccessKey"

	subject  = "SES-Mail Subject"
	bodyText = "SES-Mail body"
	bodyHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>SES-Mail</title>
</head>
<body>
SES-Mail body
</body>
</html>`
)

var toAddresses = []string{"test1@aws.com", "test2@aws.com"}

func Test_SendMail(t *testing.T) {
	api := NewAPI(endpoint, from, accessKeyID, secretAccessKey)

	body, err := api.SendMail(subject, bodyText, toAddresses)
	if err != nil {
		t.Fatalf("send text mail failed, err: %s", err.Error())
	}
	t.Logf("send text email response: %s\n", body)
}

func Test_SendHTMLMail(t *testing.T) {
	api := NewAPI(endpoint, from, accessKeyID, secretAccessKey)

	body, err := api.SendHTMLMail(subject, bodyHTML, toAddresses)
	if err != nil {
		t.Fatalf("send html mail failed, err: %s", err.Error())
	}
	t.Logf("send html email response: %s\n", body)
}
