package main

import (
	"github.com/dipperin/go-ms-toolkit/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"read-test-server/common"
	"read-test-server/router"
)

const (
	region     = "ap-east-1"
	bucketName = "hanzi-read-test"
	maxKeys    = 0

	accessKeyId     = "AKIAZBI2PYWUKHGJY5MC"
	accessSecretKey = "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"
)

func init() {
	common.InitLogger()
	common.InitAudioUploadRoot()
}

func main() {
	log := common.Log
	defer log.Sync()

	if len(os.Args) != 2 {
		log.Panic("please specify configuration file")
	}
	var conf common.ServerConfig
	confPath := os.Args[1]
	if confData, err := ioutil.ReadFile(confPath); err != nil {
		log.Panic("read config file failed", zap.Error(err))
	} else if err = json.ParseJsonFromBytes(confData, &conf); err != nil {
		log.Panic("parse config file failed", zap.Error(err))
	}

	r := gin.Default()
	router.RegisterRouter(r)
	log.Info("start gin server...")
	if err := r.Run(":1234"); err != nil {
		panic(err)
	}

	//// Load the SDK's configuration from environment and shared config, and
	//// create the client with this.
	//cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region),
	//	config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
	//		Value: aws.Credentials{
	//			AccessKeyID: accessKeyId, SecretAccessKey: accessSecretKey, SessionToken: "",
	//		},
	//	}),
	//)
	//if err != nil {
	//	log.Fatalf("failed to load SDK configuration, %v", err)
	//}
	//
	//client := s3.NewFromConfig(cfg)
	//
	//var bn = bucketName
	//// Set the parameters based on the CLI flag inputs.
	//params := &s3.ListObjectsV2Input{
	//	Bucket: &bn,
	//}
	////if len(objectPrefix) != 0 {
	////	params.Prefix = &objectPrefix
	////}
	////if len(objectDelimiter) != 0 {
	////	params.Delimiter = &objectDelimiter
	////}
	//
	//// Create the Paginator for the ListObjectsV2 operation.
	//p := s3.NewListObjectsV2Paginator(client, params, func(o *s3.ListObjectsV2PaginatorOptions) {
	//	if v := int32(maxKeys); v != 0 {
	//		o.Limit = v
	//	}
	//})
	//
	//var i int
	//for p.HasMorePages() {
	//	i++
	//	// Next Page takes a new context for each page retrieval. This is where
	//	// you could add timeouts or deadlines.
	//	page, err := p.NextPage(context.TODO())
	//	if err != nil {
	//		log.Fatalf("failed to get page %v, %v", i, err)
	//	}
	//
	//	// Log the objects found
	//	for _, obj := range page.Contents {
	//		fmt.Println("Object:", *obj.Key)
	//	}
	//}
	//
	//fileName := "test/something_to_delete.zip"
	//file, err := os.Open("./idman638build18.exe")
	//if err != nil {
	//	log.Fatalf("open file failed %v", err)
	//}
	//
	//begin := time.Now()
	//
	//object, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
	//	Key:    &fileName,
	//	Body:   file,
	//	Bucket: &bn,
	//	//ACL:                       "",
	//})
	//if err != nil {
	//	log.Fatalf("upload failed, %v", err)
	//}
	//fmt.Printf("??Object??, %#v, 耗时: %v", object, time.Now().Sub(begin))
}
