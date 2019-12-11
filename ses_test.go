package ses

import (
	"testing"
)

const (
	tEndpoint        = "endpoint"
	tFrom            = "from"
	tAccessKeyId     = "accessKeyId"
	tSecretAccessKey = "secretAccessKey"

	tSubject  = "SES-Mail Subject"
	tBodyText = "SES-Mail body"
	tBodyHtml = `<!DOCTYPE html>
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

var tToAddresses = []string{"test1@aws.com", "test2@aws.com"}

func TestSendMail(t *testing.T) {
	api := NewApi(tEndpoint,
		tFrom,
		tAccessKeyId,
		tSecretAccessKey)

	body, err := api.SendMail(tSubject, tBodyText, tToAddresses)
	if err != nil {
		t.Fatal(err)
	}

	body, err = api.SendHtmlMail(tSubject, tBodyHtml, tToAddresses)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(body)
}
