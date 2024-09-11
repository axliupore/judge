package client

import (
	"fmt"
	"github.com/axliupore/judge/pkg/cmd"
	"github.com/bytedance/sonic"
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
		client:  resty.New(),             // Initialize HTTP client
		baseURL: "http://127.0.0.1:5050", // Set base URL from configuration
	}
}

func (s *Service) Send(r *cmd.Request) (*cmd.Response, error) {

	url := fmt.Sprintf("%s%s", s.baseURL, "/run")
	rsp, err := s.client.R().SetBody(r).Post(url)
	if err != nil {
		return nil, err
	}

	res := make([]*cmd.Response, 0)
	if err = sonic.Unmarshal(rsp.Body(), &res); err != nil {
		return nil, err
	}
	return res[0], nil
}

func (s *Service) Delete(fileId string) error {

	url := fmt.Sprintf("%s/file/%s", s.baseURL, fileId)
	_, err := s.client.R().Delete(url)
	return err
}
