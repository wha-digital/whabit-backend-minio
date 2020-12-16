package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"git.innovasive.co.th/backend/helper"
	"git.innovasive.co.th/backend/minio"
)

func SetBucketPublicPolicy(client *minio.Client, bucketName string) error {
	var buf bytes.Buffer
	bu, _ := ioutil.ReadFile("./policy.json")
	t, err := template.New("policy").Parse(string(bu))
	if err != nil {
		return err
	}

	if err := t.Execute(&buf, bucketName); err != nil {
		return err
	}

	policy := buf.String()
	if err := client.GetClient().SetBucketPolicy(bucketName, policy); err != nil {
		return err
	}
	log.Println("create bucket with policy success")
	fmt.Println(policy)
	return nil
}

func main() {
	minioEndpoint := "saansook-s3-storage-dev.innovasive.in.th"
	minioAccess := "ssa-dev"
	minioSecret := "SSA#SecrEt@Inn@vAsive2o20!"
	minioSSL := true
	minioRegion := "ap-southeast-1"

	client, err := minio.NewMinio(minioEndpoint, minioAccess, minioSecret, minioSSL, minioRegion)
	if err != nil {
		panic(err)
	}
	log.Println(client)

	if err := SetBucketPublicPolicy(client, "test-policy"); err != nil {
		fmt.Println(err.Error())
		return
	}

	reader, _ := os.Open("./doctor.jpg")
	stat, _ := reader.Stat()
	size := stat.Size()
	buf, contentType, extension, err := helper.GetMimeType(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer reader.Close()
	filename := fmt.Sprintf("test%s", extension)
	fmt.Println(filename, size, contentType)

	if err := client.UploadFileWithReader("test-policy", filename, &buf, size, contentType); err != nil {
		fmt.Println(err)
		return
	}

	// if err := client.GetClient().RemoveObject("test-policy", "Insee.png"); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
