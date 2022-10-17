package router

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type progressWriter struct {
	written int64
	writer  io.WriterAt
	size    int64
}

func (pw *progressWriter) writeAt(p []byte, off int64) (int, error) {
	atomic.AddInt64(&pw.written, int64(len(p)))

	percentage_downloaded := float32(pw.written*100) / float32(pw.size)
	fmt.Printf("File size:%d downloaded:%d percentage:%.2f%%\r", pw.size, pw.written, percentage_downloaded)

	return pw.writer.WriteAt(p, off)
}

func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func getFileSize(svc *s3.S3, bucket string, prefix string) (filesize int64, error error) {
	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(prefix),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func parseFilename(key_string string) (filename string) {
	ss := strings.Split(key_string, "/")
	s := ss[len(ss)-1]
	return s
}

func Trigger(c *gin.Context) {
	bucket := os.Getenv("BUCKET")
	key := "100.json"

	filename := parseFilename(key)

	fmt.Println(os.Getenv("ACCESS_TOKEN_ID"), os.Getenv("SECRET_ACCESS_KEY"), os.Getenv("SESSION_TOKEN"))
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-3"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_TOKEN_ID"), os.Getenv("SECRET_ACCESS_KEY"), os.Getenv("SESSION_TOKEN")),
	}))

	s3_client := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)
	size, err := getFileSize(s3_client, bucket, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Starting download, size:", byteCountDecimal(size))
	cwd, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	temp, err := ioutil.TempFile(cwd, "getObjWithProgress-tmp-")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	temp_filename := temp.Name()

	writer := &progressWriter{writer: temp, size: size, written: 0}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.Download(writer.writer, params); err != nil {
		log.Printf("Download failed! Deleting tempfile: %s", temp_filename)
		os.Remove(temp_filename)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := temp.Close(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.Rename(temp.Name(), filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, filename)

}
