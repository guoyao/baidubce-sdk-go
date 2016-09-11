package bos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/util"
)

// Endpoints of baidubce
var Endpoint = map[string]string{
	"bj": "bj.bcebos.com",
	"gz": "gz.bcebos.com",
}

type Config struct {
	*bce.Config
}

func NewConfig(config *bce.Config) *Config {
	return &Config{config}
}

// Client is the client for bos.
type Client struct {
	*bce.Client
}

// DefaultClient provided a default `bos.Client` instance.
//var DefaultClient = NewClient(bce.DefaultConfig)

// NewClient returns an instance of type `bos.Client`.
func NewClient(config *Config) *Client {
	bceClient := bce.NewClient(config.Config)
	return &Client{bceClient}
}

func checkBucketName(bucketName string) {
	if bucketName == "" {
		panic("bucket name should not be empty.")
	}

	if strings.Index(bucketName, "/") == 0 {
		panic("bucket name should not be start with '/'")
	}
}

func checkObjectKey(objectKey string) {
	if objectKey == "" {
		panic("object key should not be empty.")
	}

	if strings.Index(objectKey, "/") == 0 {
		panic("object key should not be start with '/'")
	}
}

// GetBucketName returns the actual name of bucket.
func (c *Client) GetBucketName(bucketName string) string {
	if c.Endpoint != "" && !util.MapContains(Endpoint, func(key, value string) bool {
		return strings.ToLower(value) == strings.ToLower(c.Endpoint)
	}) {
		bucketName = ""
	}

	return bucketName
}

func (c *Client) GetURL(bucketName, objectKey string, params map[string]string) string {
	host := c.Endpoint

	if host == "" {
		host = Endpoint[c.GetRegion()]

		if bucketName != "" {
			host = bucketName + "." + host
		}
	}

	uriPath := objectKey

	return c.Client.GetURL(host, uriPath, params)
}

// GetBucketLocation returns the location of a bucket.
func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, *bce.Error) {
	bucketName = c.GetBucketName(bucketName)
	params := map[string]string{"location": ""}

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, "", params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var location *Location
	err = json.Unmarshal(bodyContent, &location)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return location, nil
}

// ListBuckets is for getting a collection of bucket.
func (c *Client) ListBuckets(option *bce.SignOption) (*BucketSummary, *bce.Error) {
	req, err := bce.NewRequest(http.MethodGet, c.GetURL("", "", nil), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var bucketSummary *BucketSummary
	err = json.Unmarshal(bodyContent, &bucketSummary)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return bucketSummary, nil
}

// CreateBucket is for creating a bucket.
func (c *Client) CreateBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest(http.MethodPut, c.GetURL(bucketName, "", nil), nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) DoesBucketExist(bucketName string, option *bce.SignOption) (bool, *bce.Error) {
	req, err := bce.NewRequest(http.MethodHead, c.GetURL(bucketName, "", nil), nil)

	if err != nil {
		return false, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if res != nil {
		switch {
		case res.StatusCode < http.StatusBadRequest || res.StatusCode == http.StatusForbidden:
			return true, nil
		case res.StatusCode == http.StatusNotFound:
			return false, nil
		}
	}

	return false, bceError
}

func (c *Client) DeleteBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest(http.MethodDelete, c.GetURL(bucketName, "", nil), nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) SetBucketPrivate(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["Private"], option)
}

func (c *Client) SetBucketPublicRead(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["PublicRead"], option)
}

func (c *Client) SetBucketPublicReadWrite(bucketName string, option *bce.SignOption) *bce.Error {
	return c.setBucketAclFromString(bucketName, CannedAccessControlList["PublicReadWrite"], option)
}

func (c *Client) GetBucketAcl(bucketName string, option *bce.SignOption) (*BucketAcl, *bce.Error) {
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, "", params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var bucketAcl *BucketAcl
	err = json.Unmarshal(bodyContent, &bucketAcl)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return bucketAcl, nil
}

func (c *Client) SetBucketAcl(bucketName string, bucketAcl BucketAcl, option *bce.SignOption) *bce.Error {
	byteArray, err := util.ToJson(bucketAcl, "accessControlList")

	if err != nil {
		return bce.NewError(err)
	}

	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest(http.MethodPut, c.GetURL(bucketName, "", params), bytes.NewReader(byteArray))

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) PutObject(bucketName, objectKey string, data interface{},
	metadata *ObjectMetadata, option *bce.SignOption) (PutObjectResponse, *bce.Error) {

	checkObjectKey(objectKey)

	var reader io.Reader

	if str, ok := data.(string); ok {
		reader = strings.NewReader(str)
	} else if byteArray, ok := data.([]byte); ok {
		reader = bytes.NewReader(byteArray)
	} else if r, ok := data.(io.Reader); ok {
		reader = r
	} else {
		panic("data type should be string or []byte or io.Reader.")
	}

	req, err := bce.NewRequest(http.MethodPut, c.GetURL(bucketName, objectKey, nil), reader)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if c.Checksum {
		option.AddHeader("x-bce-content-sha256", util.GetSha256(data))
	}

	if metadata != nil {
		metadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	putObjectResponse := NewPutObjectResponse(res.Header)

	return putObjectResponse, nil
}

func (c *Client) DeleteObject(bucketName, objectKey string, option *bce.SignOption) *bce.Error {
	checkObjectKey(objectKey)

	req, err := bce.NewRequest(http.MethodDelete, c.GetURL(bucketName, objectKey, nil), nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) DeleteMultipleObjects(bucketName string, objectKeys []string,
	option *bce.SignOption) (*DeleteMultipleObjectsResponse, *bce.Error) {

	checkBucketName(bucketName)

	keys := make([]string, 0, len(objectKeys))

	for _, key := range objectKeys {
		if key != "" {
			keys = append(keys, key)
		}
	}

	objectKeys = keys
	length := len(objectKeys)

	if length == 0 {
		return nil, nil
	}

	objectMap := make(map[string][]map[string]string, 1)
	objects := make([]map[string]string, length, length)

	for index, value := range objectKeys {
		objects[index] = map[string]string{"key": value}
	}

	objectMap["objects"] = objects
	byteArray, err := util.ToJson(objectMap)

	if err != nil {
		return nil, bce.NewError(err)
	}

	params := map[string]string{"delete": ""}
	body := bytes.NewReader(byteArray)

	req, err := bce.NewRequest(http.MethodPost, c.GetURL(bucketName, "", params), body)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	if len(bodyContent) > 0 {
		var deleteMultipleObjectsResponse *DeleteMultipleObjectsResponse
		err := json.Unmarshal(bodyContent, &deleteMultipleObjectsResponse)

		if err != nil {
			return nil, bce.NewError(err)
		}

		return deleteMultipleObjectsResponse, nil
	}

	return nil, nil
}

func (c *Client) ListObjects(bucketName string, option *bce.SignOption) (*ListObjectsResponse, *bce.Error) {
	return c.ListObjectsFromRequest(ListObjectsRequest{BucketName: bucketName}, option)
}

func (c *Client) ListObjectsFromRequest(listObjectsRequest ListObjectsRequest,
	option *bce.SignOption) (*ListObjectsResponse, *bce.Error) {

	bucketName := listObjectsRequest.BucketName
	params := make(map[string]string)

	if listObjectsRequest.Delimiter != "" {
		params["delimiter"] = listObjectsRequest.Delimiter
	}

	if listObjectsRequest.Marker != "" {
		params["marker"] = listObjectsRequest.Marker
	}

	if listObjectsRequest.Prefix != "" {
		params["prefix"] = listObjectsRequest.Prefix
	}

	if listObjectsRequest.MaxKeys > 0 {
		params["maxKeys"] = strconv.Itoa(listObjectsRequest.MaxKeys)
	}

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, "", params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var listObjectsResponse *ListObjectsResponse
	err = json.Unmarshal(bodyContent, &listObjectsResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return listObjectsResponse, nil
}

func (c *Client) CopyObject(srcBucketName, srcKey, destBucketName, destKey string,
	option *bce.SignOption) (*CopyObjectResponse, *bce.Error) {

	return c.CopyObjectFromRequest(CopyObjectRequest{
		SrcBucketName:  srcBucketName,
		SrcKey:         srcKey,
		DestBucketName: destBucketName,
		DestKey:        destKey,
	}, option)
}

func (c *Client) CopyObjectFromRequest(copyObjectRequest CopyObjectRequest,
	option *bce.SignOption) (*CopyObjectResponse, *bce.Error) {

	checkBucketName(copyObjectRequest.SrcBucketName)
	checkBucketName(copyObjectRequest.DestBucketName)
	checkObjectKey(copyObjectRequest.SrcKey)
	checkObjectKey(copyObjectRequest.DestKey)

	req, err := bce.NewRequest(http.MethodPut, c.GetURL(copyObjectRequest.DestBucketName, copyObjectRequest.DestKey, nil), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	source := util.URIEncodeExceptSlash(fmt.Sprintf("/%s/%s", copyObjectRequest.SrcBucketName,
		copyObjectRequest.SrcKey))

	option.AddHeader("x-bce-copy-source", source)
	copyObjectRequest.mergeToSignOption(option)

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var copyObjectResponse *CopyObjectResponse
	err = json.Unmarshal(bodyContent, &copyObjectResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return copyObjectResponse, nil
}

func (c *Client) GetObject(bucketName, objectKey string, option *bce.SignOption) (*Object, *bce.Error) {
	return c.GetObjectFromRequest(GetObjectRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}, option)
}

func (c *Client) GetObjectFromRequest(getObjectRequest GetObjectRequest,
	option *bce.SignOption) (*Object, *bce.Error) {

	checkBucketName(getObjectRequest.BucketName)
	checkObjectKey(getObjectRequest.ObjectKey)

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(getObjectRequest.BucketName, getObjectRequest.ObjectKey, nil), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	object := &Object{
		ObjectMetadata: NewObjectMetadataFromHeader(res.Header),
		ObjectContent:  res.Body,
	}

	return object, nil
}

func (c *Client) GetObjectToFile(getObjectRequest *GetObjectRequest, file *os.File,
	option *bce.SignOption) (*ObjectMetadata, *bce.Error) {

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	checkBucketName(getObjectRequest.BucketName)
	checkObjectKey(getObjectRequest.ObjectKey)

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(getObjectRequest.BucketName, getObjectRequest.ObjectKey, nil), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return objectMetadata, bce.NewError(err)
	}

	_, err = file.Write(bodyContent)

	if err != nil {
		return objectMetadata, bce.NewError(err)
	}

	return objectMetadata, nil
}

func (c *Client) GetObjectMetadata(bucketName, objectKey string, option *bce.SignOption) (*ObjectMetadata, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest(http.MethodHead, c.GetURL(bucketName, objectKey, nil), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	return objectMetadata, nil
}

func (c *Client) GeneratePresignedUrl(bucketName, objectKey string, option *bce.SignOption) (string, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, objectKey, nil), nil)

	if err != nil {
		return "", bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("host")

	authorization := bce.GenerateAuthorization(*c.Credentials, *req, option)
	url := fmt.Sprintf("%s?authorization=%s", req.URL.String(), util.URLEncode(authorization))

	return url, nil
}

func (c *Client) AppendObject(bucketName, objectKey string, offset int, data interface{},
	metadata *ObjectMetadata, option *bce.SignOption) (AppendObjectResponse, *bce.Error) {

	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	var reader io.Reader

	if str, ok := data.(string); ok {
		reader = strings.NewReader(str)
	} else if byteArray, ok := data.([]byte); ok {
		reader = bytes.NewReader(byteArray)
	} else if r, ok := data.(io.Reader); ok {
		byteArray, err := ioutil.ReadAll(r)

		if err != nil {
			return nil, bce.NewError(err)
		}

		reader = bytes.NewReader(byteArray)
	} else {
		panic("data type should be string or []byte or io.Reader.")
	}

	params := map[string]string{"append": ""}

	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}

	req, err := bce.NewRequest(http.MethodPost, c.GetURL(bucketName, objectKey, params), reader)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if c.Checksum {
		option.AddHeader("x-bce-content-sha256", util.GetSha256(data))
	}

	if metadata != nil {
		metadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	appendObjectResponse := NewAppendObjectResponse(res.Header)

	return appendObjectResponse, nil
}

func (c *Client) InitiateMultipartUpload(initiateMultipartUploadRequest InitiateMultipartUploadRequest,
	option *bce.SignOption) (*InitiateMultipartUploadResponse, *bce.Error) {

	bucketName := initiateMultipartUploadRequest.BucketName
	objectKey := initiateMultipartUploadRequest.ObjectKey

	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	params := map[string]string{"uploads": ""}

	req, err := bce.NewRequest(http.MethodPost, c.GetURL(bucketName, objectKey, params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if initiateMultipartUploadRequest.ObjectMetadata != nil {
		initiateMultipartUploadRequest.ObjectMetadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var initiateMultipartUploadResponse *InitiateMultipartUploadResponse
	err = json.Unmarshal(bodyContent, &initiateMultipartUploadResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return initiateMultipartUploadResponse, nil
}

func (c *Client) UploadPart(uploadPartRequest UploadPartRequest,
	option *bce.SignOption) (UploadPartResponse, *bce.Error) {

	bucketName := uploadPartRequest.BucketName
	objectKey := uploadPartRequest.ObjectKey
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	if uploadPartRequest.PartNumber < MIN_PART_NUMBER || uploadPartRequest.PartNumber > MAX_PART_NUMBER {
		panic(fmt.Sprintf("Invalid partNumber %d. The valid range is from %d to %d.",
			uploadPartRequest.PartNumber, MIN_PART_NUMBER, MAX_PART_NUMBER))
	}

	if uploadPartRequest.PartSize > 1024*1024*1024*5 {
		panic(fmt.Sprintf("PartNumber %d: Part Size should not be more than 5GB.", uploadPartRequest.PartSize))
	}

	params := map[string]string{
		"partNumber": strconv.Itoa(uploadPartRequest.PartNumber),
		"uploadId":   uploadPartRequest.UploadId,
	}

	req, err := bce.NewRequest(http.MethodPut, c.GetURL(bucketName, objectKey, params), uploadPartRequest.PartData)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeaders(map[string]string{
		"Content-Length": strconv.FormatInt(uploadPartRequest.PartSize, 10),
		"Content-Type":   "application/octet-stream",
	})

	if _, ok := option.Headers["Content-MD5"]; !ok {
		option.AddHeader("Content-MD5", util.GetMD5(uploadPartRequest.PartData, true))
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	uploadPartResponse := NewUploadPartResponse(res.Header)

	return uploadPartResponse, nil
}

func (c *Client) CompleteMultipartUpload(completeMultipartUploadRequest CompleteMultipartUploadRequest,
	option *bce.SignOption) (*CompleteMultipartUploadResponse, *bce.Error) {

	bucketName := completeMultipartUploadRequest.BucketName
	objectKey := completeMultipartUploadRequest.ObjectKey
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	completeMultipartUploadRequest.sort()
	params := map[string]string{"uploadId": completeMultipartUploadRequest.UploadId}
	byteArray, err := util.ToJson(completeMultipartUploadRequest, "parts")

	if err != nil {
		return nil, bce.NewError(err)
	}

	req, err := bce.NewRequest(http.MethodPost, c.GetURL(bucketName, objectKey, params), bytes.NewReader(byteArray))

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var completeMultipartUploadResponse *CompleteMultipartUploadResponse

	err = json.Unmarshal(bodyContent, &completeMultipartUploadResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return completeMultipartUploadResponse, nil
}

func (c *Client) MultipartUploadFromFile(bucketName, objectKey, filePath string,
	partSize int64) (*CompleteMultipartUploadResponse, *bce.Error) {

	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
	}

	initiateMultipartUploadResponse, bceError := c.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

	if bceError != nil {
		return nil, bceError
	}

	uploadId := initiateMultipartUploadResponse.UploadId

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return nil, bce.NewError(err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var totalSize int64 = fileInfo.Size()
	var partCount int = int(math.Ceil(float64(totalSize) / float64(partSize)))

	parts := make([]PartSummary, 0, partCount)

	var waitGroup sync.WaitGroup

	for i := 0; i < partCount; i++ {
		var skipBytes int64 = partSize * int64(i)
		var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

		tempFile, err := util.TempFile(nil, "", "")

		if err != nil {
			return nil, bce.NewError(err)
		}

		limitReader := io.LimitReader(file, size)
		_, err = io.Copy(tempFile, limitReader)

		if err != nil {
			return nil, bce.NewError(err)
		}

		partNumber := i + 1

		uploadPartRequest := UploadPartRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
			PartSize:   size,
			PartNumber: partNumber,
			PartData:   tempFile,
		}

		waitGroup.Add(1)

		parts = append(parts, PartSummary{PartNumber: partNumber})

		go func(partNumber int, f *os.File) {
			defer func() {
				f.Close()
				os.Remove(f.Name())
				waitGroup.Done()
			}()

			uploadPartResponse, bceError := c.UploadPart(uploadPartRequest, nil)
			uploadPartRequest.PartData = nil

			if bceError != nil {
				panic(bceError)
			}

			parts[partNumber-1].ETag = uploadPartResponse.GetETag()
		}(partNumber, tempFile)
	}

	waitGroup.Wait()
	waitGroup.Add(1)

	var completeMultipartUploadResponse *CompleteMultipartUploadResponse

	go func() {
		defer waitGroup.Done()

		completeMultipartUploadRequest := CompleteMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
			Parts:      parts,
		}

		completeResponse, bceError := c.CompleteMultipartUpload(completeMultipartUploadRequest, nil)

		if bceError != nil {
			panic(bceError)
		}

		completeMultipartUploadResponse = completeResponse
	}()

	waitGroup.Wait()

	return completeMultipartUploadResponse, nil
}

func (c *Client) AbortMultipartUpload(abortMultipartUploadRequest AbortMultipartUploadRequest,
	option *bce.SignOption) *bce.Error {

	bucketName := abortMultipartUploadRequest.BucketName
	objectKey := abortMultipartUploadRequest.ObjectKey
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	params := map[string]string{"uploadId": abortMultipartUploadRequest.UploadId}

	req, err := bce.NewRequest(http.MethodDelete, c.GetURL(bucketName, objectKey, params), nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option)

	return bceError
}

func (c *Client) ListParts(bucketName, objectKey, uploadId string,
	option *bce.SignOption) (*ListPartsResponse, *bce.Error) {

	return c.ListPartsFromRequest(ListPartsRequest{
		BucketName: bucketName,
		ObjectKey:  objectKey,
		UploadId:   uploadId,
	}, option)
}

func (c *Client) ListPartsFromRequest(listPartsRequest ListPartsRequest,
	option *bce.SignOption) (*ListPartsResponse, *bce.Error) {

	bucketName := listPartsRequest.BucketName
	objectKey := listPartsRequest.ObjectKey

	params := map[string]string{"uploadId": listPartsRequest.UploadId}

	if listPartsRequest.PartNumberMarker != "" {
		params["partNumberMarker"] = listPartsRequest.PartNumberMarker
	}

	if listPartsRequest.MaxParts > 0 {
		params["maxParts"] = strconv.Itoa(listPartsRequest.MaxParts)
	}

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, objectKey, params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var listPartsResponse *ListPartsResponse

	err = json.Unmarshal(bodyContent, &listPartsResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return listPartsResponse, nil
}

func (c *Client) ListMultipartUploads(bucketName string,
	option *bce.SignOption) (*ListMultipartUploadsResponse, *bce.Error) {

	return c.ListMultipartUploadsFromRequest(ListMultipartUploadsRequest{BucketName: bucketName}, option)
}

func (c *Client) ListMultipartUploadsFromRequest(listMultipartUploadsRequest ListMultipartUploadsRequest,
	option *bce.SignOption) (*ListMultipartUploadsResponse, *bce.Error) {

	bucketName := listMultipartUploadsRequest.BucketName

	params := map[string]string{"uploads": ""}

	if listMultipartUploadsRequest.Delimiter != "" {
		params["delimiter"] = listMultipartUploadsRequest.Delimiter
	}

	if listMultipartUploadsRequest.KeyMarker != "" {
		params["keyMarker"] = listMultipartUploadsRequest.KeyMarker
	}

	if listMultipartUploadsRequest.Prefix != "" {
		params["prefix"] = listMultipartUploadsRequest.Prefix
	}

	if listMultipartUploadsRequest.MaxUploads > 0 {
		params["maxUploads"] = strconv.Itoa(listMultipartUploadsRequest.MaxUploads)
	}

	req, err := bce.NewRequest(http.MethodGet, c.GetURL(bucketName, "", params), nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option)

	if bceError != nil {
		return nil, bceError
	}

	bodyContent, err := res.GetBodyContent()

	if err != nil {
		return nil, bce.NewError(err)
	}

	var listMultipartUploadsResponse *ListMultipartUploadsResponse

	err = json.Unmarshal(bodyContent, &listMultipartUploadsResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return listMultipartUploadsResponse, nil
}

func (c *Client) setBucketAclFromString(bucketName, acl string, option *bce.SignOption) *bce.Error {
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest(http.MethodPut, c.GetURL(bucketName, "", params), nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	headers := map[string]string{"x-bce-acl": acl}
	option.AddHeaders(headers)

	_, bceError := c.SendRequest(req, option)

	return bceError
}
