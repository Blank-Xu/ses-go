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

func TestSendMail(t *testing.T) {
	api := NewAPI(endpoint, from, accessKeyID, secretAccessKey)

	body, err := api.SendMail(subject, bodyText, toAddresses)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(body)

	body2, err := api.SendHTMLMail(subject, bodyHTML, toAddresses)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(body2)
}
