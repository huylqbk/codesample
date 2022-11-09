package minio

import (
	"bytes"
	"fmt"
	"log"

	"github.com/minio/minio-go"
)

type service struct {
	client *minio.Client
	config *Config
}

type Config struct {
	Addr          string `mapstructure:"addr"`
	AccessKey     string `mapstructure:"accessKey"`
	SecretKey     string `mapstructure:"secretKey"`
	UseSSL        bool   `mapstructure:"useSSL"`
	DefaultBucket string `mapstructure:"defaultBucket"`
}

func NewService(cfg *Config) *service {
	endpoint := cfg.Addr
	accessKeyID := cfg.AccessKey
	secretAccessKey := cfg.SecretKey
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, cfg.UseSSL)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := &service{
		client: minioClient,
		config: cfg,
	}
	if cfg.DefaultBucket != "" {
		if !s.BucketExist(cfg.DefaultBucket) {
			s.CreateBucket(cfg.DefaultBucket)
		}
	}
	return s
}

func (s *service) Get(bucket, path string) ([]byte, error) {
	obj, err := s.client.GetObject(bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(obj)
	return buf.Bytes(), nil
}

func (s *service) Put(bucket, path string, data []byte) error {
	file := bytes.NewReader(data)
	info, err := s.client.PutObject(bucket, path, file, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	log.Println("Put file success", info)
	return nil
}

func (s *service) Remove(bucket, path string) error {
	return s.client.RemoveObject(bucket, path)
}

func (s *service) GetFile(path string) ([]byte, error) {
	obj, err := s.client.GetObject(s.config.DefaultBucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(obj)
	return buf.Bytes(), nil
}

func (s *service) PutFile(path string, data []byte) error {
	file := bytes.NewReader(data)
	info, err := s.client.PutObject(s.config.DefaultBucket, path, file, int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	log.Println("Put file success", info)
	return nil
}

func (s *service) RemoveFile(path string) error {
	return s.client.RemoveObject(s.config.DefaultBucket, path)
}

func (s *service) CreateBucket(name string) error {
	err := s.client.MakeBucket(name, "us-east-1")
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("Create bucket success", name)
	return nil
}

func (s *service) BucketExist(name string) bool {
	found, err := s.client.BucketExists(name)
	if err != nil {
		log.Println(err)
		return false
	}
	return found
}
