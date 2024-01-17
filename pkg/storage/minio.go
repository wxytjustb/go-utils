package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	once           sync.Once
	minioClient    *MinioClient
	minioConfig    *MinioConfig
	bucketCheckMap = make(map[string]struct{})
)

type MinioConfig struct {
	Endpoints      string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKey      string `mapstructure:"access-key" json:"access-key" yaml:"access-key"`
	SecretKey      string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	UseSSL         bool   `mapstructure:"useSSL" json:"useSSL" yaml:"useSSL" default:"true"`
	ObjectLocking  bool   `mapstructure:"object-locking" json:"object-locking" yaml:"object-locking" default:"false"`
	TimeoutSeconds int    `mapstructure:"timeout" json:"timeout" yaml:"timeout" default:"60"`
}

func MinioModuleInitialize(config *MinioConfig) (err error) {
	minioConfig = config

	once.Do(func() {
		minioClient, err = NewMinioClient()
	})

	return err
}

type MinioClient struct {
	client *minio.Client
}

func GetSingleMinioClient() *MinioClient {
	return minioClient
}

func NewMinioClient() (*MinioClient, error) {
	client, err := minio.New(minioConfig.Endpoints, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: minioConfig.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return &MinioClient{client: client}, nil
}

// EnableBucket 检查bucket是否存在，不存在就创建
func (m *MinioClient) EnableBucket(ctx context.Context, bucketName string) (err error) {

	// skip check if bucket is exist
	if _, ok := bucketCheckMap[bucketName]; ok {
		return nil
	}

	isExist, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !isExist {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region:        "",
			ObjectLocking: minioConfig.ObjectLocking,
		})
		if err != nil {
			return err
		}
	}

	bucketCheckMap[bucketName] = struct{}{}

	return err
}

// PutObject 上传文件
func (m *MinioClient) PutObject(ctx context.Context, path string, reader *io.Reader) (err error) {

	bucketName, objectName, err := m.GetObjectInfoByPath(path)
	if err != nil {
		return err
	}

	_, err = m.client.PutObject(ctx, bucketName, objectName, *reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	return err
}

// GetPreSignedPath 获取文件的预签名地址,也就是临时访问路径
func (m *MinioClient) GetPreSignedPath(ctx context.Context, path string, duration time.Duration) (string, error) {

	paths := strings.SplitN(path, "/", 3)
	if len(paths) != 3 {
		return "", fmt.Errorf("invalid path")
	}

	url, err := m.client.PresignedGetObject(ctx, paths[1], paths[2], duration, nil)
	if err != nil {
		fmt.Println(err)
	}

	return fmt.Sprintf("%s?%s", url.Path, url.RawFragment), nil
}

// GetObjectInfoByPath 通过路径获取bucketName和objectName
func (m *MinioClient) GetObjectInfoByPath(path string) (bucketName, objectName string, err error) {
	paths := strings.SplitN(path, "/", 3)
	if len(paths) != 3 {
		return "", "", fmt.Errorf("invalid path")
	}
	return paths[1], paths[2], nil
}
