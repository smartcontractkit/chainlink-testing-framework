package s3provider

import (
	"context"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMinioFactory(t *testing.T) {
	port := freeport.GetOne(t)
	consolePort := freeport.GetOne(t)
	s3provider, err := NewMinioFactory().NewProvider(WithPort(port), WithConsolePort(consolePort))
	require.NoError(t, err)

	t.Logf("URL: %s", s3provider.GetURL())

	// Initialize minio client object.
	minioClient, err := minio.New(s3provider.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(s3provider.GetAccessKey(), s3provider.GetSecretKey(), ""),
		Secure: false,
	})
	require.NoError(t, err)

	// Test file upload
	filename := "test.txt"
	filePath := "./" + filename
	contentType := "application/octet-stream"

	info, err := minioClient.FPutObject(
		context.Background(),
		s3provider.GetBucket(),
		filename,
		filePath,
		minio.PutObjectOptions{ContentType: contentType},
	)
	require.NoError(t, err)
	require.Equal(t, int64(7), info.Size)

	t.Logf("successfully uploaded %s of size %d bytes\n", filename, info.Size)
}
