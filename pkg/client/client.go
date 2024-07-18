package client

import (
	"encoding/json"
	"fmt"
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/pkg/exec"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/run"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
)

// Service defines the client service structure.
type Service struct {
	client  *resty.Client // HTTP client instance
	baseURL string        // Base URL for API requests
}

// NewService creates a new Service instance.
func NewService() *Service {
	return &Service{
		client:  resty.New(),                     // Initialize HTTP client
		baseURL: config.CoreConfig.Judge.Address, // Set base URL from configuration
	}
}

// sendRequest sends a request and handles the response.
func (s *Service) sendRequest(endpoint string, body interface{}, response interface{}) error {
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint) // Construct full URL for the request

	rsp, err := s.client.R().
		SetBody(body).
		Post(url) // Send POST request with body
	if err != nil {
		return err
	}

	err = json.Unmarshal(rsp.Body(), response) // Parse response body into provided response struct
	if err != nil {
		return err
	}
	return nil
}

// SendRunRequest sends a run request to the judge service.
func (s *Service) SendRunRequest(req *run.Request) (*run.Response, error) {
	var res []run.Response
	log.Logger.Infof("sending run request: %s", spew.Sdump(req))
	err := s.sendRequest("/run", req, &res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	log.Logger.Infof("accept run response: %s", spew.Sdump(res))
	return &res[0], nil
}

// SendExecRequest sends a test execution request to the judge service.
func (s *Service) SendExecRequest(req *exec.Request) (*exec.Response, error) {
	var res []exec.Response
	log.Logger.Infof("sending exec request: %s", spew.Sdump(req))
	err := s.sendRequest("/run", req, &res)
	if err != nil || len(res) == 0 {
		return &exec.Response{}, err
	}
	log.Logger.Infof("accept exec response: %s", spew.Sdump(res))
	return &res[0], nil
}

// DeleteRequest sends a file deletion request to the judge service.
func (s *Service) DeleteRequest(fileId string) error {
	url := fmt.Sprintf("%s/file/%s", s.baseURL, fileId) // Construct URL for file deletion
	_, err := s.client.R().Delete(url)                  // Send DELETE request
	return err
}
