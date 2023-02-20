package api

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
)

type SonarLogin struct {
	BaseUrl string
	Token   string
}

func (s *SonarLogin) Get(router string, v url.Values) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	apiUrl := strings.Join([]string{s.BaseUrl, router}, "/")
	req, err := http.NewRequest("GET", apiUrl, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.Token, "")
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (s *SonarLogin) Post(router string, formData url.Values) (*http.Response, error) {
	client := &http.Client{}
	apiUrl := strings.Join([]string{s.BaseUrl, router}, "/")
	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.Token, "")
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
