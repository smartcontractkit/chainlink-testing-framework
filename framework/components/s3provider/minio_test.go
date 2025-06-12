package s3provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	port := freeport.GetOne(t)
	consolePort := freeport.GetOne(t)
	accessKey, secretKey := randomStr(accessKeyLength), randomStr(secretKeyLength)
	s3provider, err := NewMinioFactory().New(
		WithPort(port),
		WithConsolePort(consolePort),
		WithAccessKey(accessKey),
		WithSecretKey(secretKey),
	)
	require.NoError(t, err)

	// Test Output
	output := s3provider.Output()
	expected := &Output{
		AccessKey:      accessKey,
		SecretKey:      secretKey,
		Bucket:         DefaultBucket,
		ConsoleURL:     s3provider.GetConsoleURL(),
		ConsoleBaseURL: s3provider.GetConsoleBaseURL(),
		Endpoint:       s3provider.GetEndpoint(),
		BaseEndpoint:   fmt.Sprintf("%s:%d", DefaultHost, port),
		Region:         s3provider.GetRegion(),
		UseCache:       false,
	}
	fmt.Printf("%#v\n%#v\n", expected, output)
	require.True(t, cmp.Equal(expected, output))
	require.Len(t, output.AccessKey, accessKeyLength)
	require.Len(t, output.SecretKey, secretKeyLength)

	// Initialize minio client object.
	minioClient, err := minio.New(s3provider.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(s3provider.GetAccessKey(), s3provider.GetSecretKey(), ""),
		Secure: false,
	})
	require.NoError(t, err)

	info, err := helperUploadFile(minioClient, s3provider.GetBucket())
	require.NoError(t, err)
	require.Equal(t, int64(7), info.Size)

	statusCode, err := helperDownloadFile(
		fmt.Sprintf(
			"http://%s/%s/%s",
			output.Endpoint,
			output.Bucket,
			info.Key,
		),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
}

func TestNewFrom(t *testing.T) {
	port := freeport.GetOne(t)
	consolePort := freeport.GetOne(t)

	input := &Input{
		Port:        port,
		ConsolePort: consolePort,
	}

	output, err := NewMinioFactory().NewFrom(input)
	require.NoError(t, err)

	// Test Output
	expected := &Output{
		Bucket:         DefaultBucket,
		ConsoleURL:     fmt.Sprintf("http://%s:%d", "127.0.0.1", consolePort),
		ConsoleBaseURL: fmt.Sprintf("http://%s:%d", DefaultHost, consolePort),
		Endpoint:       fmt.Sprintf("%s:%d", "127.0.0.1", port),
		BaseEndpoint:   fmt.Sprintf("%s:%d", DefaultHost, port),
		Region:         DefaultRegion,
		UseCache:       false,
	}
	fmt.Printf("%#v\n%#v\n", expected, output)
	require.True(t, cmp.Equal(expected, output, cmpopts.IgnoreFields(Output{}, "AccessKey", "SecretKey")))
	require.Len(t, output.AccessKey, accessKeyLength)
	require.Len(t, output.SecretKey, secretKeyLength)

	// Initialize minio client object.
	minioClient, err := minio.New(output.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(output.AccessKey, output.SecretKey, ""),
		Secure: false,
	})
	require.NoError(t, err)

	info, err := helperUploadFile(minioClient, output.Bucket)
	require.NoError(t, err)
	require.Equal(t, int64(7), info.Size)

	statusCode, err := helperDownloadFile(
		fmt.Sprintf(
			"http://%s/%s/%s",
			output.Endpoint,
			output.Bucket,
			info.Key,
		),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, statusCode)
}

func helperUploadFile(minioClient *minio.Client, bucket string) (*minio.UploadInfo, error) {
	// Test file upload
	filename := "test.txt"
	filePath := "./" + filename
	contentType := "application/octet-stream"
	info, err := minioClient.FPutObject(
		context.Background(),
		bucket,
		filename,
		filePath,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func helperDownloadFile(url string) (int, error) {
	fmt.Printf("Downloading: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return resp.StatusCode, nil
}
