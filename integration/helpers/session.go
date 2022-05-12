package helpers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Session struct {
	host    string
	headers http.Header
}

func NewSession(host string) Session {
	return Session{host: host}
}

func (s *Session) Login(user, password string) error {
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
	header := http.Header{}
	header.Add(http.CanonicalHeaderKey("Authorization"), auth)
	s.headers = header
	return nil
}

func (s *Session) Logout() {
	s.headers = http.Header{}
}

func (s *Session) Get(path string) (int, io.Reader, error) {
	parsedURL, err := url.Parse(s.host + path)
	if err != nil {
		return 0, nil, err
	}

	var c http.Client
	resp, err := c.Do(&http.Request{
		Method: http.MethodGet,
		URL:    parsedURL,
		Header: s.headers,
	})
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}

func (s *Session) Delete(path string) (int, error) {
	parsedURL, err := url.Parse(s.host + path)
	if err != nil {
		return 0, err
	}

	var c http.Client
	resp, err := c.Do(&http.Request{
		Method: http.MethodDelete,
		URL:    parsedURL,
		Header: s.headers,
	})
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func (s *Session) Post(path string, payload interface{}) (int, io.Reader, error) {
	parsedURL, err := url.Parse(s.host + path)
	if err != nil {
		return 0, nil, err
	}

	// Some magic so the function have the right type.
	// It's done at a lower lever with the http.Client for a POST
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(payload)
	bufNC := io.NopCloser(buf)

	var c http.Client
	s.headers.Set(http.CanonicalHeaderKey("Content-Type"), "application/json")
	defer s.headers.Del(http.CanonicalHeaderKey("Content-Type"))

	resp, err := c.Do(&http.Request{
		Method: http.MethodPost,
		URL:    parsedURL,
		Body:   bufNC,
		Header: s.headers,
	})

	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, resp.Body, nil
}

func (s *Session) PostMultipart(path, title, filename string, buff io.Reader) (int, io.Reader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Title
	if err := writer.WriteField("title", title); err != nil {
		return 0, nil, err
	}

	// Video
	fileWriter, _ := writer.CreateFormFile("video", filename)
	if _, err := io.Copy(fileWriter, buff); err != nil {
		return 0, nil, err
	}
	if err := writer.Close(); err != nil {
		return 0, nil, err
	}

	contentType := fmt.Sprintf("multipart/form-data; boundary=%s", writer.Boundary())
	r, err := http.NewRequest(http.MethodPost, s.host+path, bytes.NewReader(body.Bytes()))
	if err != nil {
		return 0, nil, err
	}
	s.headers.Set(http.CanonicalHeaderKey("Content-Type"), contentType)
	defer s.headers.Del(http.CanonicalHeaderKey("Content-Type"))
	r.Header = s.headers

	var c http.Client
	resp, err := c.Do(r)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, resp.Body, nil
}

func (s *Session) WaitVideoEncoded(path string) error {
	var videoStatus VideoStatus
	for strings.ToLower(videoStatus.Status) != "complete" {
		time.Sleep(5 * time.Second)
		code, body, err := s.Get(path)
		if err != nil {
			return err
		}

		if code != 200 {
			return fmt.Errorf("wrong response, status code = %d", code)
		}

		// Reading the body
		rawBody, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rawBody, &videoStatus)
		if err != nil {
			return err
		}
	}
	return nil
}
