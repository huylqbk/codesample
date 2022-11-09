package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type service struct {
	client *elasticsearch.Client
	config *Config
}

type Config struct {
	Addresses     []string      `mapstructure:"addresses" validate:"required"`
	Username      string        `mapstructure:"username"`
	Password      string        `mapstructure:"password"`
	APIKey        string        `mapstructure:"apiKey"`
	Header        http.Header   // Global HTTP request header.
	EnableLogging bool          `mapstructure:"enableLogging"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

type MultiMatchQuery struct {
	Query Query `json:"query"`
}

type Query struct {
	Bool Bool `json:"bool"`
}

type Bool struct {
	Must []any `json:"must"`
}

type MultiMatch struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
}

type EsHits[T any] struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewService(cfg *Config) (*service, error) {
	config := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		APIKey:    cfg.APIKey,
		Header:    cfg.Header,
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Second * 5
	}

	if cfg.EnableLogging {
		config.Logger = &elastictransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &service{
		client: client,
		config: cfg,
	}, nil
}

func (s *service) Info(ctx context.Context) (*esapi.Response, error) {
	response, err := s.client.Info(s.client.Info.WithContext(ctx), s.client.Info.WithHuman())
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, errors.New(response.String())
	}

	return response, nil
}

func (s *service) Index(ctx context.Context, name string, data []byte) (*esapi.Response, error) {
	response, err := s.client.Indices.Create(
		name,
		s.client.Indices.Create.WithContext(ctx),
		s.client.Indices.Create.WithBody(bytes.NewReader(data)),
		s.client.Indices.Create.WithPretty(),
		s.client.Indices.Create.WithHuman(),
		s.client.Indices.Create.WithTimeout(s.config.Timeout),
	)
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, errors.New(response.String())
	}

	return response, nil
}

func (s *service) Alias(ctx context.Context, indexes []string, name string, data []byte) (*esapi.Response, error) {
	response, err := s.client.Indices.PutAlias(
		indexes,
		name,
		s.client.Indices.PutAlias.WithBody(bytes.NewReader(data)),
		s.client.Indices.PutAlias.WithContext(ctx),
		s.client.Indices.PutAlias.WithHuman(),
		s.client.Indices.PutAlias.WithPretty(),
		s.client.Indices.PutAlias.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, errors.New(response.String())
	}

	return response, nil
}

func (s *service) Search(ctx context.Context, index, term string, fields []string) (*esapi.Response, error) {
	query := MultiMatchQuery{
		Query: Query{
			Bool: Bool{
				Must: []any{MultiMatch{
					Query:  term,
					Fields: fields,
				}},
			},
		},
	}

	dataBytes, err := json.Marshal(&query)
	if err != nil {
		return nil, err
	}

	response, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(index),
		s.client.Search.WithBody(bytes.NewReader(dataBytes)),
		s.client.Search.WithPretty(),
		s.client.Search.WithHuman(),
		s.client.Search.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	if response.IsError() {
		return nil, errors.Wrap(errors.New(response.String()), "s.client.Search error")
	}

	return response, nil
}

func (s *service) Exsist(ctx context.Context, indexes []string) (*esapi.Response, error) {
	response, err := s.client.Indices.Exists(
		indexes,
		s.client.Indices.Exists.WithContext(ctx),
		s.client.Indices.Exists.WithHuman(),
		s.client.Indices.Exists.WithPretty(),
	)
	if err != nil {
		return nil, err
	}

	if response.IsError() && response.StatusCode != 404 {
		return nil, errors.New(response.String())
	}

	return response, nil
}
