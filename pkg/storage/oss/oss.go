package oss

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yeisme/notevault/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// minioClient Global MinIO client
var minioClient *minio.Client

// InitOss Initialize MinIO object storage client
func InitOss(ossConfig config.OssConfig) error {
	// Parse endpoint to extract host and protocol
	u, err := url.Parse(ossConfig.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint: %v", err)
	}

	host := u.Host
	useSSL := strings.HasPrefix(ossConfig.Endpoint, "https://")

	// Create MinIO client
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(ossConfig.AccessKeyID, ossConfig.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %v", err)
	}

	// Save to global variable
	minioClient = client

	// Check if bucket exists, create it if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, ossConfig.BucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket status: %v", err)
	}

	if !exists {
		logx.Infof("bucket %s does not exist, creating...", ossConfig.BucketName)
		err = client.MakeBucket(ctx, ossConfig.BucketName, minio.MakeBucketOptions{
			Region: ossConfig.Region,
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		logx.Infof("bucket %s created successfully", ossConfig.BucketName)
	} else {
		logx.Infof("bucket %s already exists", ossConfig.BucketName)
	}

	return nil
}

// GetOssClient Returns the MinIO client, if not initialized, returns nil
func GetOssClient() *minio.Client {
	if minioClient == nil {
		logx.Error("MinIO client is not initialized")
		return nil
	}
	return minioClient
}
