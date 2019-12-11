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

func NewApi(endpoint, from, accessKeyId, secretAccessKey string) *Option {
	op := new(Option)
	op.endpoint = endpoint
	op.source = from
	op.accessKeyId = accessKeyId
	op.secretAccessKey = []byte(secretAccessKey)
	return op
}

type Option struct {
	endpoint        string
	source          string
	accessKeyId     string
	secretAccessKey []byte
}

func (p *Option) SendMail(subject, bodyText string, toAddresses []string) (string, error) {
	subject = fmt.Sprintf("=?UTF-8?B?%s?=",
		base64.StdEncoding.EncodeToString([]byte(subject)))

	data := make(url.Values, 5+len(toAddresses))
	data.Add("Action", "SendEmail")
	data.Add("Source", p.source)
	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Text.Data", bodyText)
	data.Add("AWSAccessKeyId", p.accessKeyId)

	for idx, email := range toAddresses {
		data.Add("Destination.ToAddresses.member."+strconv.Itoa(idx+1), email)
	}

	body, err := p.doPost(data)

	return string(body), err
}

func (p *Option) SendHtmlMail(subject, bodyHtml string, toAddresses []string) (string, error) {
	subject = fmt.Sprintf("=?UTF-8?B?%s?=",
		base64.StdEncoding.EncodeToString([]byte(subject)))

	data := make(url.Values, 5+len(toAddresses))
	data.Add("Action", "SendEmail")
	data.Add("Source", p.source)
	data.Add("Message.Subject.Data", subject)
	data.Add("Message.Body.Html.Data", bodyHtml)
	data.Add("AWSAccessKeyId", string(p.accessKeyId))

	for idx, email := range toAddresses {
		data.Add("Destination.ToAddresses.member."+strconv.Itoa(idx+1), email)
	}

	body, err := p.doPost(data)

	return string(body), err
}

func (p *Option) doGet(data url.Values) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		p.endpoint,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	return p.doRequest(req)
}

func (p *Option) doPost(data url.Values) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		p.endpoint,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	return p.doRequest(req)
}

var (
	defaultClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * 30,
	}
)

func (p *Option) doRequest(request *http.Request) ([]byte, error) {
	// set Header Content-Type
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// set Header Date
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 -0700")
	request.Header.Set("Date", date)

	// calculate signature
	h := hmac.New(sha256.New, p.secretAccessKey)
	h.Write([]byte(date))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// calculate Authorization
	// fmt.Sprintf("AWS3-HTTPS AWSAccessKeyId=%s, Algorithm=HmacSHA256, Signature=%s", "", "")
	var authorizationBuf bytes.Buffer
	authorizationBuf.Grow(128)
	authorizationBuf.WriteString("AWS3-HTTPS AWSAccessKeyId=")
	authorizationBuf.WriteString(p.accessKeyId)
	authorizationBuf.WriteString(", Algorithm=HmacSHA256, Signature=")
	authorizationBuf.WriteString(signature)

	// set Header X-Amzn-Authorization
	request.Header.Set("X-Amzn-Authorization", authorizationBuf.String())

	resp, err := defaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request failed, err: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || err != nil {
		return nil, fmt.Errorf("http response failed, code: %d, body: %s, err: %v", resp.StatusCode, body, err)
	}

	return body, nil
}
