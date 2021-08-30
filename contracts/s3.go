package contracts

import (
	"archive/zip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// S3Downloader S3 artifacts syncing for external contract sources
type S3Downloader struct {
	Region     string
	Session    *session.Session
	Downloader *s3manager.Downloader
	Bucket     string
}

// NewS3Downloader creates new S3 downloader
func NewS3Downloader(cfg config.ExternalSources) *S3Downloader {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region)},
	)
	return &S3Downloader{
		Session:    sess,
		Downloader: s3manager.NewDownloader(sess),
		Bucket:     cfg.S3URL,
	}
}

// UpdateSources downloads contracts artifacts by commit
func (d *S3Downloader) UpdateSources(cfg config.ExternalSources) error {
	for repoRootDir, s := range cfg.Repositories {
		artifactName := d.artifactWithCommit(repoRootDir, s.Commit)
		fp := filepath.Join(cfg.RootPath, repoRootDir, artifactName)
		log.Debug().Str("Archive", fp).Msg("Downloading artifact")
		file, err := os.Create(fp)
		if err != nil {
			return err
		}
		_, err = d.Downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(d.Bucket),
				Key:    aws.String(artifactName),
			})
		if err != nil {
			return errors.Wrapf(err, "artifact %s not found in bucket %s", artifactName, d.Bucket)
		}
		dir, _ := filepath.Split(fp)
		_, err = Unzip(fp, dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *S3Downloader) artifactWithCommit(name string, commit string) string {
	return fmt.Sprintf("%s-%s.zip", name, commit)
}

// Unzip decompress a zip archive to dest
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	//nolint
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		if _, err = io.Copy(outFile, rc); err != nil {
			return filenames, err
		}
		if err = outFile.Close(); err != nil {
			return filenames, err
		}
		if err = rc.Close(); err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
