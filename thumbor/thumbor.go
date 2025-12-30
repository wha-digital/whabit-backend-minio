package thumbor

import (
	"encoding/json"
	"regexp"
	"strings"
)

type Thumbor struct {
	MinioEndPoint string
	EndPoint      string
	SecurityKey   string
	ImgWidth      string
	ImgHeight     string
	FilterConfig  string
}

func NewThumbor(minioEndpoint string, minioSSL string, thumborEndpoint string, securityKey string, imgWidth string, imgHeight string, filterConfig string) *Thumbor {
	var minioURL string
	if minioSSL == "true" {
		minioURL = "https://"
	} else {
		minioURL = "http://"
	}

	minioURL = minioURL + minioEndpoint
	return &Thumbor{
		MinioEndPoint: minioURL,
		EndPoint:      thumborEndpoint,
		SecurityKey:   securityKey,
		ImgWidth:      imgWidth,
		ImgHeight:     imgHeight,
		FilterConfig:  filterConfig,
	}
}

func (t *Thumbor) NewImageLink(bucketName string, uri string) string {
	/* if svg file return miniopath */
	if uri != "" {
		ex := regexp.MustCompile(`.svg`)
		if ex.MatchString(uri) {
			return t.MinioEndPoint + "/" + bucketName + "/" + uri
		}
	}

	if uri != "" {
		ex := regexp.MustCompile(`.pdf`)
		if ex.MatchString(uri) {
			return t.MinioEndPoint + "/" + bucketName + "/" + uri
		}
	}

	if uri != "" {
		if bucketName != "" {
			if !strings.Contains(uri, "://") {
				uri = bucketName + "/" + uri
			}

		}
		if string([]rune(uri)[0]) != "/" && !strings.Contains(uri, "://") {
			uri = "/" + uri
		}
		if !strings.Contains(uri, "://") {
			uri = t.MinioEndPoint + uri
		}
	}
	link := t.EndPoint + "/" + t.SecurityKey + "/"
	if t.ImgWidth != "" || t.ImgHeight != "" {
		link = link + t.ImgWidth + "x" + t.ImgHeight
	}
	link = link + t.FilterConfig + "/" + uri
	return link
}

func (t Thumbor) String() string {
	bu, _ := json.Marshal(t)
	return string(bu)
}
