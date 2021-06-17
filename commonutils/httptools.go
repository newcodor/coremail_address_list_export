package commonutils

import (
	"bufio"
	"crypto/tls"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	httpClient        *http.Client
	globalHttpTimeout int = 7
)

func init() {
	httpClient, _ = NewHttpClient(7, "")
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	b, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(b, "")
	return e
}
func NewHttpClient(httpTimeout int, proxy string) (*http.Client, error) {
	if httpTimeout == 0 {
		httpTimeout = globalHttpTimeout
	}
	// var proxyUrl *url.URL = nil
	var proxyUrl *url.URL = nil
	var err error = nil
	if proxy != "" {
		proxyUrl, _ = url.Parse(proxy)
	}
	return &http.Client{
		Timeout: time.Duration(httpTimeout) * time.Second,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}, err
}

func HttpPost(url string, params string, headers map[string]string, timeout int, httpclient *http.Client) (map[string]string, error) {
	respList := make(map[string]string)
	if httpclient != nil {
		httpClient = httpclient
	}
	if timeout > 0 {
		httpClient.Timeout = time.Duration(timeout) * time.Second

	}
	req, reqErr := http.NewRequest(http.MethodPost, url, strings.NewReader(params))
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if headers != nil {
		for header, headerValue := range headers {
			req.Header.Set(header, headerValue)
		}
	}
	resp, err := httpClient.Do(req)
	if reqErr == nil && err == nil {
		respList["url"] = url
		respList["statusCode"] = strconv.Itoa(resp.StatusCode)
		reader := bufio.NewReader(resp.Body)
		e := determineEncoding(reader)
		uft8Reader := transform.NewReader(reader, e.NewDecoder())
		respBody, _ := ioutil.ReadAll(uft8Reader)
		respList["respBody"] = string(respBody)
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	return respList, err
}
