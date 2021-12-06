package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var TIMEOUT = 30

type Auth struct {
	Username string
	Password string
}

func NewAuth(username, password string) *Auth {
	return &Auth{
		Username: username,
		Password: password,
	}
}

func HTTPGet(uri string, params map[string]string, headers map[string]string) (*http.Response, error) {
	return HTTPGetBasicAuth(uri, params, headers, nil)
}

func HTTPGetBasicAuth(uri string, params map[string]string, headers map[string]string, auth *Auth) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	q := req.URL.Query()

	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	// Set the auth for the request.
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	client := NewClient(time.Duration(TIMEOUT) * time.Second)
	return client.Do(req)
}

func HTTPPost(url string, body interface{}, params map[string]string, headers map[string]string, client *http.Client) (*http.Response, error) {
	var bodyJSON []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJSON, err = json.Marshal(body)
		if err != nil {
			log.Println(err)
			return nil, errors.New("http post body to json failed")
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJSON))

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf(fmt.Sprintf("new request is fail: %v \n", err))
	}

	req.Header.Set("Content-type", "application/json;charset=utf-8")
	q := req.URL.Query()

	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	return client.Do(req)
}

func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				err = conn.SetDeadline(time.Now().Add(timeout))
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
			ResponseHeaderTimeout: timeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout * 2,
	}
}

func newRequest(url string, body interface{}, params map[string]string, headers map[string]string, method string) (*http.Request, error) {
	var read io.Reader
	var req *http.Request
	if body != nil {
		switch body.(type) {
		case multipart.File:
			file, err := ioutil.ReadAll(body.(multipart.File))
			if err != nil {
				return nil, err
			}
			read = bytes.NewReader(file)
		case []byte:
			read = bytes.NewReader(body.([]byte))
		default:
			bodyJSON, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			read = bytes.NewBuffer(bodyJSON)
		}
	}

	req, err := http.NewRequest(method, url, read)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}

	q := req.URL.Query()

	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	return req, nil
}

func HTTPBasicAuthorization(userName string, password string, url string, client *http.Client) bool {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Println(err)
		return false
	}

	client.Jar = jar

	req, err := newRequest(url, nil, nil, nil, "GET")
	if err != nil {
		log.Println(err)
		return false
	}
	req.SetBasicAuth(userName, password)
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
