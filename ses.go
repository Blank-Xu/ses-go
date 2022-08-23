package ses

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NewAPI  constructor for SES api
func NewAPI(endpoint, from, accessKeyID, secretAccessKey string) *Option {
	op := new(Option)

	op.endpoint = endpoint
	op.source = from
	op.accessKeyID = accessKeyID
	op.secretAccessKey = []byte(secretAccessKey)

	return op
}

// NewHTTPClientAPI  constructor for SES api with given HTTP client
func NewHTTPClientAPI(client http.Client, endpoint, from, accessKeyID, secretAccessKey string) *Option {
	SetDefaultHTTPClient(client)

	return NewAPI(endpoint, from, accessKeyID, secretAccessKey)
}

// SetDefaultHTTPClient  set the default http client for this package
func SetDefaultHTTPClient(client http.Client) {
	defaultHTTPClient = client
}

var defaultHTTPClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: time.Second * 30,
}

// Option  SES api options
type Option struct {
	endpoint        string
	source          string
	accessKeyID     string
	secretAccessKey []byte
}

// SendMail  send text emails
func (p *Option) SendMail(subject, bodyText string, toAddresses []string) (string, error) {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(subject)))
	base64.StdEncoding.Encode(b, []byte(subject))

	var subBuf bytes.Buffer

	subBuf.Grow(128)
	subBuf.WriteString("=?UTF-8?B?")
	subBuf.Write(b)
	subBuf.WriteString("?=")

	data := make(url.Values, 6+len(toAddresses))
	data.Add("Action", "SendEmail")
	data.Add("Source", p.source)
	data.Add("Message.Subject.Data", subBuf.String())
	data.Add("Message.Body.Text.Data", bodyText)
	data.Add("AWSAccessKeyId", p.accessKeyID)

	for idx, email := range toAddresses {
		data.Add("Destination.ToAddresses.member."+strconv.Itoa(idx+1), email)
	}

	body, err := p.doRequest(data)

	return string(body), err
}

// SendHTMLMail  send html emails
func (p *Option) SendHTMLMail(subject, bodyHTML string, toAddresses []string) (string, error) {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(subject)))
	base64.StdEncoding.Encode(b, []byte(subject))

	var subBuf bytes.Buffer

	subBuf.Grow(128)
	subBuf.WriteString("=?UTF-8?B?")
	subBuf.Write(b)
	subBuf.WriteString("?=")

	data := make(url.Values, 6+len(toAddresses))
	data.Add("Action", "SendEmail")
	data.Add("Source", p.source)
	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Html.Data", bodyHTML)
	data.Add("AWSAccessKeyId", string(p.accessKeyID))

	for idx, email := range toAddresses {
		data.Add("Destination.ToAddresses.member."+strconv.Itoa(idx+1), email)
	}

	body, err := p.doRequest(data)

	return string(body), err
}

func (p *Option) doRequest(data url.Values) ([]byte, error) {
	request, err := http.NewRequest(
		http.MethodPost,
		p.endpoint,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	// set Header Content-Type
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// set Header Date
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 -0700")
	request.Header.Set("Date", date)

	// calculate signature
	h := hmac.New(sha256.New, p.secretAccessKey)
	_, _ = h.Write([]byte(date))

	signature := make([]byte, base64.StdEncoding.EncodedLen(h.Size()))
	base64.StdEncoding.Encode(signature, h.Sum(nil))

	// calculate Authorization
	var authorizationBuf bytes.Buffer

	authorizationBuf.Grow(128)
	authorizationBuf.WriteString("AWS3-HTTPS AWSAccessKeyId=")
	authorizationBuf.WriteString(p.accessKeyID)
	authorizationBuf.WriteString(", Algorithm=HmacSHA256, Signature=")
	authorizationBuf.Write(signature)

	// set Header X-Amzn-Authorization
	request.Header.Set("X-Amzn-Authorization", authorizationBuf.String())

	resp, err := defaultHTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request failed, err: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || err != nil {
		return nil, fmt.Errorf("http response failed, code: %d, body: %s, err: %v", resp.StatusCode, body, err)
	}

	return body, nil
}
