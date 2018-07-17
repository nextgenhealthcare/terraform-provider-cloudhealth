package cloudhealth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// AwsExternalID is used to enable integration with AWS via IAM Roles.
type AwsExternalID struct {
	ExternalID string `json:"generated_external_id"`
}

// GetAwsExternalID gets the AWS External ID tied to the CloudHealth Account.
func (s *Client) GetAwsExternalID() (string, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("aws_accounts/:id/generate_external_id?api_key=%s", s.ApiKey))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", url.String(), nil)

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var id = new(AwsExternalID)
		err = json.Unmarshal(responseBody, &id)
		if err != nil {
			return "", err
		}

		return id.ExternalID, nil
	case http.StatusUnauthorized:
		return "", ErrClientAuthenticationError
	case http.StatusForbidden:
		return "", ErrClientAuthenticationError
	default:
		return "", fmt.Errorf("Unknown Response with CloudHealth: `%d`", resp.StatusCode)
	}
}
