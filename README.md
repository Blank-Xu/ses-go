# SES

[![Go Report Card](https://goreportcard.com/badge/github.com/Blank-Xu/ses-go)](https://goreportcard.com/report/github.com/Blank-Xu/ses-go)
[![PkgGoDev](https://pkg.go.dev/badge/Blank-Xu/ses-go)](https://pkg.go.dev/github.com/Blank-Xu/ses-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

---

An easy way to use the [Amazon Simple Email Service(SES)](https://aws.amazon.com/ses/) api to send emails.

## Installation

    go get github.com/Blank-Xu/ses-go
    
## Simple Example
```go
package main

import (
	ses "github.com/Blank-Xu/ses-go"
	
	"fmt"
)

func main() {
	api := ses.NewAPI(endpoint, from, accessKeyID, secretAccessKey)

	body, err := api.SendMail(subject, bodyText, toAddresses)
	fmt.Println(body, err)

	body2, err := api.SendHTMLMail(subject, bodyHTML, toAddresses)
	fmt.Println(body2, err)
}
```

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.