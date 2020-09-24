package minio

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"mime/multipart"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	helperModel "git.innovasive.co.th/backend/models"
	"github.com/minio/minio-go"
	minioLib "github.com/minio/minio-go"
)

const (
	MINIO_DEFAULT_REGION = "ap-southeast-1"
)

type Client struct {
	MinioClient    *minioLib.Client
	MinioEndPoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioSSL       bool
	Region         string
}

func NewMinio(endpoint string, access string, secret string, ssl bool, region string) (*Client, error) {

	min, err := minioLib.NewWithRegion(endpoint, access, secret, ssl, region)
	if err != nil {
		return nil, err
	}

	c := &Client{
		MinioClient:    min,
		MinioEndPoint:  endpoint,
		MinioAccessKey: access,
		MinioSecretKey: secret,
		MinioSSL:       ssl,
		Region:         region,
	}

	return c, nil
}

func (c *Client) GetMinioURI() string {
	var minioEndURI string
	if c.MinioSSL {
		minioEndURI = "https://"
	} else {
		minioEndURI = "http://"
	}
	minioEndURI = minioEndURI + c.MinioEndPoint
	return minioEndURI
}

func CreateImageFile(img image.Image, format string, filename string, desc string) error {
	file, err := os.Create(desc)
	if err != nil {
		return err
	}
	switch format {
	case "jpg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
		break
	case "jpeg":
		if err := jpeg.Encode(file, img, nil); err != nil {
			return err
		}
		break
	case "png":
		if err := png.Encode(file, img); err != nil {
			return err
		}
		break
	}
	defer file.Close()

	return nil
}

/*
	This function will generate path filename for upload minio
	format
		YYYYMMDD_id_[0-9]{10}_extension

		- YYYY ปีแบบคริสศักราช
		- MM เดือน
		- DD วันที่
		- id string
		- [0-9]{10} เลข generate 0-9 จำนวน 10 ตัว
		- extension
	example
		20200413-9026362617-example.jpg
*/
func generateObjectName(foldername string, id string, extension string) string {
	var date = strings.ReplaceAll(helperModel.NewDateFromTime(time.Now()).String(), "-", "")
	var generateNumber = strconv.Itoa(rand.Intn(1000000000))
	foldername = func() string {
		lastString := []rune(foldername)
		if string(lastString) == "/" {
			return foldername
		}
		return foldername + "/"
	}()

	return fmt.Sprintf(foldername+"%s_%s_%s.%s", date, id, generateNumber, extension)
}

func GenerateObjectName(foldername string, id string, filename string) string {
	return generateObjectName(foldername, id, filename)
}

func (c *Client) GenerateObjectName(foldername string, id string, filename string) string {
	return generateObjectName(foldername, id, filename)
}

func (c *Client) GetClient() *minioLib.Client {
	return c.MinioClient
}

func (c *Client) GetEndPoint() string {
	return c.MinioEndPoint
}

func (c *Client) UploadMultipartFile(bucketName string, objectName string, file *multipart.FileHeader) (err error) {
	contentType := file.Header.Get("Content-Type")
	size := file.Size

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err = c.GetClient().PutObject(bucketName, objectName, src, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return err
	}
	return nil
}
func (c *Client) UploadFromFile(bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(bucketName, objectName, pathFile, minio.PutObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *Client) UploadFromFilePDF(bucketName string, foldername string, pathFile string, filename string) error {
	objectName := foldername + "/" + filename

	src, err := os.Open(pathFile)
	if err != nil {
		return err
	}

	defer src.Close()

	if _, err := c.GetClient().FPutObject(bucketName, objectName, pathFile, minio.PutObjectOptions{ContentType: "application/pdf", ContentEncoding: "UTF-8"}); err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateBucket(bucketName string, region string) error {
	if region == "" {
		region = MINIO_DEFAULT_REGION
	}
	if err := c.GetClient().MakeBucket(bucketName, region); err != nil {
		return err
	}
	return nil
}

func (c *Client) ExistBucket(bucketName string) (bool, error) {
	exists, err := c.GetClient().BucketExists(bucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Client) GetPathImageFromURL(url string) string {
	var pathImageRegex = regexp.MustCompile(c.MinioEndPoint + `\/[a-z\-]+\/`)
	var path string
	if url != "" {
		if pathImageRegex.MatchString(url) {
			loc := pathImageRegex.FindStringIndex(url)
			if len(loc) > 0 {
				path = url[loc[1]:]
			}
		}
	}
	return path
}
