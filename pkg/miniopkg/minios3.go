package miniopkg

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOCreds struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type MinIOObjOptions struct {
	MinClient   *minio.Client
	Location    string
	Bucketname  string
	ObjectName  string
	Filepath    string
	ContentType string
}

// NewMinIOClient creates a new client to interact with MINIO
func (mc *MinIOCreds) NewMinIOClient() *minio.Client {
	minioClient, err := minio.New(mc.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(mc.AccessKeyID, mc.SecretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("unable to create MINIO Client: ", err)
	}
	return minioClient
}

// UploadObject uploads the object to Minio ObjectStore
func (ob *MinIOObjOptions) UploadObject() {
	ctx := context.Background()
	//Check if bucket exists
	bucket, err := ob.MinClient.BucketExists(ctx, ob.Bucketname)
	if err != nil {
		panic(err)
	}
	if !bucket {
		err := ob.MinClient.MakeBucket(ctx, ob.Bucketname, minio.MakeBucketOptions{Region: ob.Location})
		if err != nil {
			fmt.Println(err)
		}
	}
	obj, err := ob.MinClient.FPutObject(ctx, ob.Bucketname, ob.ObjectName, ob.Filepath,
		minio.PutObjectOptions{ContentType: ob.ContentType})
	if err != nil {
		panic(err)
	}
	fmt.Println("OBJECT KEY: ", obj.Key)
	fmt.Printf("Successfully Uploaded '%s' to the bucket '%s' of size '%d' Bytes\n", obj.Key, obj.Bucket, obj.Size)
}

func (ob *MinIOObjOptions) GetObjects() {

	//ctx = context.Background()

}
