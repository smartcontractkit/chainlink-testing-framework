package s3provider

import (
	"context"
	"fmt"
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
	require.True(t,
		cmp.Equal(&Output{
			AccessKey:    accessKey,
			SecretKey:    secretKey,
			Bucket:       DefaultBucket,
			ConsoleURL:   s3provider.GetConsoleURL(),
			Endpoint:     s3provider.GetEndpoint(),
			BaseEndpoint: fmt.Sprintf("%s:%d", DefaultHost, port),
			Region:       s3provider.GetRegion(),
			UseCache:     false,
		}, output))
	require.Len(t, output.AccessKey, accessKeyLength)
	require.Len(t, output.SecretKey, secretKeyLength)

	// Initialize minio client object.
	minioClient, err := minio.New(s3provider.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(s3provider.GetAccessKey(), s3provider.GetSecretKey(), ""),
		Secure: false,
	})
	require.NoError(t, err)

	helperUploadFile(t, minioClient, s3provider.GetBucket())
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
	fmt.Printf("%#v\n", output)
	require.True(t,
		cmp.Equal(&Output{
			Bucket:       DefaultBucket,
			ConsoleURL:   fmt.Sprintf("http://%s:%d", "127.0.0.1", consolePort),
			Endpoint:     fmt.Sprintf("%s:%d", "127.0.0.1", port),
			BaseEndpoint: fmt.Sprintf("%s:%d", "minio", port),
			Region:       DefaultRegion,
			UseCache:     false,
		}, output, cmpopts.IgnoreFields(Output{}, "AccessKey", "SecretKey")))
	require.Len(t, output.AccessKey, accessKeyLength)
	require.Len(t, output.SecretKey, secretKeyLength)

	// Initialize minio client object.
	minioClient, err := minio.New(output.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(output.AccessKey, output.SecretKey, ""),
		Secure: false,
	})
	require.NoError(t, err)

	helperUploadFile(t, minioClient, output.Bucket)
}

func helperUploadFile(t *testing.T, minioClient *minio.Client, bucket string) {
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
	require.NoError(t, err)
	require.Equal(t, int64(7), info.Size)
}
