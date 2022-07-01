package clients

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TransformerInfos struct {
	Name    string
	Address string
	Port    string
}

type IConsulClient interface {
	GetTransformationServices() ([]TransformerInfos, error)
	GetTransformationService(name string) (TransformerInfos, error)
}

var _ IConsulClient = &consulClient{}

type consulClient struct {
	GetTransformerServices *http.Request
	GetTransformerService  *http.Request
}

func NewConsulClient(host, user, password string) (IConsulClient, error) {
	// Create header
	header := http.Header{}
	header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+password)))

	// Create GetTransformerServices url
	parsedURLServices, err := url.Parse("http://" + host + `/v1/agent/services?filter="transformer"+in+Tags`)
	if err != nil {
		return nil, err
	}

	// Create GetTransformerService url
	parsedURLService, err := url.Parse("http://" + host + `/v1/agent/services?filter="transformer"+in+Tags+and+`)
	if err != nil {
		return nil, err
	}

	// Create GetTransformerServices request
	consulC := &consulClient{
		GetTransformerServices: &http.Request{
			Method: http.MethodGet,
			URL:    parsedURLServices,
			Header: header,
		},
		GetTransformerService: &http.Request{
			Method: http.MethodGet,
			URL:    parsedURLService,
			Header: header,
		},
	}

	return consulC, nil
}

// Get all available instances of transformation services
func (c *consulClient) GetTransformationServices() ([]TransformerInfos, error) {
	var httpC http.Client
	resp, err := httpC.Do(c.GetTransformerServices)
	if err != nil {
		return nil, err
	}

	var parsedResp map[string]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return nil, err
	}

	var transformers []TransformerInfos
	for _, value := range parsedResp {
		transformers = append(transformers, TransformerInfos{
			Name:    fmt.Sprintf("%v", value["Service"]),
			Address: fmt.Sprintf("%v", value["Address"]),
			Port:    fmt.Sprintf("%v", value["Port"]),
		})
	}

	return transformers, nil
}

// Get all available instances of a given transformation service
func (c *consulClient) GetTransformationService(name string) (TransformerInfos, error) {
	// Create new request
	request, err := c.createGetServiceRequest(name)
	if err != nil {
		return TransformerInfos{}, err
	}

	var httpC http.Client
	resp, err := httpC.Do(request)
	if err != nil {
		return TransformerInfos{}, err
	}

	var parsedResp map[string]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&parsedResp)
	if err != nil {
		return TransformerInfos{}, err
	}

	var transformer TransformerInfos
	for _, value := range parsedResp {
		transformer = TransformerInfos{
			Name:    fmt.Sprintf("%v", value["Service"]),
			Address: fmt.Sprintf("%v", value["Address"]),
			Port:    fmt.Sprintf("%v", value["Port"]),
		}
	}
	return transformer, nil
}

func (c *consulClient) createGetServiceRequest(name string) (*http.Request, error) {
	// Create URL with new query parameter
	urlStr := strings.Clone(c.GetTransformerService.URL.String())
	GetServiceURL, err := url.Parse(urlStr + `"` + name + `"` + "+in+Service")
	if err != nil {
		return nil, err
	}

	// Create new request
	return &http.Request{
		Method: c.GetTransformerService.Method,
		URL:    GetServiceURL,
		Header: c.GetTransformerService.Header,
	}, nil
}
