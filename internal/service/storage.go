package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

const R2FreeTierLimit = 10 * 1024 * 1024 * 1024 // 10 GB

type FileInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	URL          string    `json:"url"`
}

type StorageInfo struct {
	UsedSpace      int64   `json:"usedSpace"`
	RemainingSpace int64   `json:"remainingSpace"`
	TotalSpace     int64   `json:"totalSpace"`
	PercentageUsed float64 `json:"percentageUsed"`
	FileCount      int     `json:"fileCount"`
}

type StorageService struct {
	client     *s3.Client
	cfg        *config.Config
	bucketName string
	baseUrl    string
}

func NewStorageService(cfg *config.Config) *StorageService {
	if !cfg.Storage.Enabled {
		return &StorageService{cfg: cfg}
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.Storage.AccountID)
	baseUrl := cfg.Storage.PublicURL
	if baseUrl == "" {
		baseUrl = fmt.Sprintf("https://%s.%s.r2.cloudflarestorage.com", cfg.Storage.Bucket, cfg.Storage.AccountID)
	}

	cfgSdk, err := s3Config.LoadDefaultConfig(context.TODO(),
		s3Config.WithRegion("auto"),
		s3Config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Storage.AccessKey,
			cfg.Storage.SecretKey,
			"",
		)),
	)
	if err != nil {
		slog.Error("Failed to initialize AWS SDK for R2", "error", err)
		panic(err)
	}

	client := s3.NewFromConfig(cfgSdk, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &StorageService{
		client:     client,
		cfg:        cfg,
		bucketName: cfg.Storage.Bucket,
		baseUrl:    baseUrl,
	}
}

func (s *StorageService) CheckConnection(ctx context.Context) (bool, error) {
	if !s.cfg.Storage.Enabled {
		return false, fmt.Errorf("R2 storage is disabled")
	}

	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		slog.Error("R2 connection check failed", "error", err)
		return false, err
	}

	return true, nil
}

func (s *StorageService) ListFiles(ctx context.Context, prefix *string) ([]FileInfo, error) {
	if !s.cfg.Storage.Enabled {
		return nil, fmt.Errorf("R2 storage is disabled")
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	}
	if prefix != nil && *prefix != "" {
		input.Prefix = prefix
	}

	resp, err := s.client.ListObjectsV2(ctx, input)
	if err != nil {
		slog.Error("Failed to list files from R2", "error", err)
		return nil, err
	}

	var files []FileInfo
	for _, item := range resp.Contents {
		files = append(files, FileInfo{
			Key:          aws.ToString(item.Key),
			Size:         aws.ToInt64(item.Size),
			LastModified: aws.ToTime(item.LastModified),
			URL:          fmt.Sprintf("%s/%s", s.baseUrl, aws.ToString(item.Key)),
		})
	}

	return files, nil
}

func (s *StorageService) GetStorageInfo(ctx context.Context) (*StorageInfo, error) {
	files, err := s.ListFiles(ctx, nil)
	if err != nil {
		return nil, err
	}

	var usedSpace int64 = 0
	for _, f := range files {
		usedSpace += f.Size
	}

	remainingSpace := R2FreeTierLimit - usedSpace
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	percentageUsed := (float64(usedSpace) / float64(R2FreeTierLimit)) * 100.0

	return &StorageInfo{
		UsedSpace:      usedSpace,
		RemainingSpace: remainingSpace,
		TotalSpace:     R2FreeTierLimit,
		PercentageUsed: math.Round(percentageUsed*100) / 100,
		FileCount:      len(files),
	}, nil
}

func (s *StorageService) UploadFile(ctx context.Context, filename string, data []byte, contentType string) (*FileInfo, error) {
	if !s.cfg.Storage.Enabled {
		return nil, fmt.Errorf("R2 storage is disabled")
	}

	info, err := s.GetStorageInfo(ctx)
	if err != nil {
		return nil, err
	}

	fileSize := int64(len(data))
	if info.UsedSpace+fileSize > R2FreeTierLimit {
		return nil, fmt.Errorf("storage is full, remaining space: %s", s.formatBytes(info.RemainingSpace))
	}

	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.ToLower(base)

	reg := regexp.MustCompile("[^a-z0-9]")
	baseClean := reg.ReplaceAllString(base, "-")

	uniqueID := uuid.New().String()[:8]
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	key := fmt.Sprintf("uploads/%d-%s-%s%s", timestamp, uniqueID, baseClean, ext)

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		slog.Error("Failed to upload file to R2", "key", key, "error", err)
		return nil, err
	}

	slog.Info("File uploaded successfully to R2", "key", key, "size", s.formatBytes(fileSize))

	return &FileInfo{
		Key:          key,
		Size:         fileSize,
		LastModified: time.Now(),
		URL:          fmt.Sprintf("%s/%s", s.baseUrl, key),
	}, nil
}

func (s *StorageService) DeleteFile(ctx context.Context, key string) error {
	if !s.cfg.Storage.Enabled {
		return fmt.Errorf("R2 storage is disabled")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		slog.Error("Failed to delete file from R2", "key", key, "error", err)
		return err
	}

	slog.Info("File deleted successfully from R2", "key", key)
	return nil
}

func (s *StorageService) GetFileInfo(ctx context.Context, key string) (*FileInfo, error) {
	if !s.cfg.Storage.Enabled {
		return nil, fmt.Errorf("R2 storage is disabled")
	}

	resp, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Key:          key,
		Size:         aws.ToInt64(resp.ContentLength),
		LastModified: aws.ToTime(resp.LastModified),
		URL:          fmt.Sprintf("%s/%s", s.baseUrl, key),
	}, nil
}

func (s *StorageService) formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 Bytes"
	}
	units := []string{"Bytes", "KB", "MB", "GB", "TB"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(1024)))
	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(1024, float64(i)), units[i])
}
