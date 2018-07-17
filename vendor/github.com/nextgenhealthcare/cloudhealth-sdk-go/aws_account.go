package cloudhealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// AwsAccount represents the configuration of an AWS Account enabled in CloudHealth.
type AwsAccount struct {
	ID             int                      `json:"id"`
	Name           string                   `json:"name"`
	Authentication AwsAccountAuthentication `json:"authentication"`
}

// AwsAccountAuthentication represents the authentication details for AWS integration.
type AwsAccountAuthentication struct {
	Protocol             string `json:"protocol"`
	AccessKey            string `json:"access_key,omitempty"`
	SecreyKey            string `json:"secret_key,omitempty"`
	AssumeRoleArn        string `json:"assume_role_arn,omitempty"`
	AssumeRoleExternalID string `json:"assume_role_external_id,omitempty"`
}

// ErrAwsAccountNotFound is returned when an AWS Account doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrAwsAccountNotFound = errors.New("AWS Account not found")

// GetAwsAccount gets the AWS Account with the specified CloudHealth Account ID.
func (s *Client) GetAwsAccount(id int) (*AwsAccount, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var account = new(AwsAccount)
		err = json.Unmarshal(responseBody, &account)
		if err != nil {
			return nil, err
		}

		return account, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusNotFound:
		return nil, ErrAwsAccountNotFound
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// CreateAwsAccount enables a new AWS Account in CloudHealth.
func (s *Client) CreateAwsAccount(account AwsAccount) (*AwsAccount, error) {

	body, _ := json.Marshal(account)

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts?api_key=%s", s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		var account = new(AwsAccount)
		err = json.Unmarshal(responseBody, &account)
		if err != nil {
			return nil, err
		}

		return account, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("Bad Request. Please check if a AWS Account with this name `%s` already exists", account.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// UpdateAwsAccount updates an existing AWS Account in CloudHealth.
func (s *Client) UpdateAwsAccount(account AwsAccount) (*AwsAccount, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", account.ID, s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	body, _ := json.Marshal(account)

	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer((body)))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var account = new(AwsAccount)
		err = json.Unmarshal(responseBody, &account)
		if err != nil {
			return nil, err
		}

		return account, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusUnprocessableEntity:
		return nil, fmt.Errorf("Bad Request. Please check if a AWS Account with this name `%s` already exists", account.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}

// DeleteAwsAccount removes the AWS Account with the specified CloudHealth ID.
func (s *Client) DeleteAwsAccount(id int) error {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/%d?api_key=%s", id, s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("DELETE", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrAwsAccountNotFound
	case http.StatusUnauthorized:
		return ErrClientAuthenticationError
	default:
		return fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}
