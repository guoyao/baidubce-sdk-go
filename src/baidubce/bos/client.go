package bos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	bce "baidubce"
	"baidubce/util"
)

// Client is the client for bos.
type Client struct {
	bce.Client
}

// DefaultClient provided a default `bos.Client` instance.
var DefaultClient = Client{bce.Client{bce.DefaultConfig}}

// NewClient returns an instance of type `bos.Client`.
func NewClient(config bce.Config) Client {
	bceClient := bce.Client{config}
	return Client{bceClient}
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
	if c.Endpoint != "" && !util.MapContains(bce.Region, func(key, value string) bool {
		return strings.ToLower(value) == strings.ToLower(c.Endpoint)
	}) {
		bucketName = ""
	}

	return bucketName
}

func (c *Client) GetBucketEndpoint(bucketName string) string {
	endpoint := c.Endpoint

	if endpoint == "" {
		endpoint = bce.Region["bj"]
	}

	if bucketName != "" {
		endpoint = bucketName + "." + endpoint
	}

	return endpoint
}

// GetBucketLocation returns the location of a bucket.
func (c *Client) GetBucketLocation(bucketName string, option *bce.SignOption) (*Location, *bce.Error) {
	bucketName = c.GetBucketName(bucketName)
	params := map[string]string{"location": ""}

	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var location *Location
	err = json.Unmarshal(res.Body, &location)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return location, nil
}

// ListBuckets is for getting a collection of bucket.
func (c *Client) ListBuckets(option *bce.SignOption) (*BucketSummary, *bce.Error) {
	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.Endpoint, nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var bucketSummary *BucketSummary
	err = json.Unmarshal(res.Body, &bucketSummary)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return bucketSummary, nil
}

// CreateBucket is for creating a bucket.
func (c *Client) CreateBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}

func (c *Client) DoesBucketExist(bucketName string, option *bce.SignOption) (bool, *bce.Error) {
	req, err := bce.NewRequest("HEAD", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return false, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if res != nil {
		switch {
		case res.StatusCode < 400 || res.StatusCode == 403:
			return true, nil
		case res.StatusCode == 404:
			return false, nil
		}
	}

	return false, bceError
}

func (c *Client) DeleteBucket(bucketName string, option *bce.SignOption) *bce.Error {
	req, err := bce.NewRequest("DELETE", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option, true)

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
	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var bucketAcl *BucketAcl
	err = json.Unmarshal(res.Body, &bucketAcl)

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
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, bytes.NewReader(byteArray))

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	_, bceError := c.SendRequest(req, option, true)

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
		byteArray, err := ioutil.ReadAll(r)

		if err != nil {
			return nil, bce.NewError(err)
		}

		reader = bytes.NewReader(byteArray)
	} else {
		panic("data type should be string or []byte or io.Reader.")
	}

	req, err := bce.NewRequest("PUT", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, reader)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if metadata != nil {
		metadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	putObjectResponse := NewPutObjectResponse(res.Header)

	return putObjectResponse, nil
}

func (c *Client) DeleteObject(bucketName, objectKey string, option *bce.SignOption) *bce.Error {
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("DELETE", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return bceError
	}

	return nil
}

func (c *Client) DeleteMultipleObjects(bucketName string, objectKeys []string,
	option *bce.SignOption) (*DeleteMultipleObjectsResponse, *bce.Error) {

	checkBucketName(bucketName)

	length := len(objectKeys)
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

	req, err := bce.NewRequest("POST", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, body)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	if len(res.Body) > 0 {
		var deleteMultipleObjectsResponse *DeleteMultipleObjectsResponse
		err := json.Unmarshal(res.Body, &deleteMultipleObjectsResponse)

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

	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var listObjectsResponse *ListObjectsResponse
	err = json.Unmarshal(res.Body, &listObjectsResponse)

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

	req, err := bce.NewRequest("PUT", c.GetUriPath(copyObjectRequest.DestKey),
		c.GetBucketEndpoint(copyObjectRequest.DestBucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	source := util.URIEncodeExceptSlash(fmt.Sprintf("/%s/%s", copyObjectRequest.SrcBucketName,
		copyObjectRequest.SrcKey))

	option.AddHeader("x-bce-copy-source", source)
	copyObjectRequest.mergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var copyObjectResponse *CopyObjectResponse
	err = json.Unmarshal(res.Body, &copyObjectResponse)

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

	req, err := bce.NewRequest(
		"GET",
		c.GetUriPath(getObjectRequest.ObjectKey),
		c.GetBucketEndpoint(getObjectRequest.BucketName),
		nil,
		nil,
	)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, false)

	if bceError != nil {
		return nil, bceError
	}

	object := &Object{
		ObjectMetadata: NewObjectMetadataFromHeader(res.Header),
		ObjectContent:  res.Response.Body,
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

	req, err := bce.NewRequest(
		"GET",
		c.GetUriPath(getObjectRequest.ObjectKey),
		c.GetBucketEndpoint(getObjectRequest.BucketName),
		nil,
		nil,
	)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	getObjectRequest.MergeToSignOption(option)

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	_, err = file.Write(res.Body)

	if err != nil {
		return objectMetadata, bce.NewError(err)
	}

	return objectMetadata, nil
}

func (c *Client) GetObjectMetadata(bucketName, objectKey string, option *bce.SignOption) (*ObjectMetadata, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("HEAD", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, false)

	if bceError != nil {
		return nil, bceError
	}

	objectMetadata := NewObjectMetadataFromHeader(res.Header)

	return objectMetadata, nil
}

func (c *Client) GeneratePresignedUrl(bucketName, objectKey string, option *bce.SignOption) (string, *bce.Error) {
	checkBucketName(bucketName)
	checkObjectKey(objectKey)

	req, err := bce.NewRequest("GET", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), nil, nil)

	if err != nil {
		return "", bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("host")

	authorization := bce.GenerateAuthorization(c.Credentials, *req, option)
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

	req, err := bce.NewRequest("POST", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), params, reader)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if metadata != nil {
		metadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option, true)

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

	req, err := bce.NewRequest("POST", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	option.AddHeader("Content-Type", util.GuessMimeType(objectKey))

	if initiateMultipartUploadRequest.ObjectMetadata != nil {
		initiateMultipartUploadRequest.ObjectMetadata.mergeToSignOption(option)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var initiateMultipartUploadResponse *InitiateMultipartUploadResponse
	err = json.Unmarshal(res.Body, &initiateMultipartUploadResponse)

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

	req, err := bce.NewRequest("PUT", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName),
		params, bytes.NewReader(uploadPartRequest.PartData))

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

	res, bceError := c.SendRequest(req, option, false)

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

	params := map[string]string{"uploadId": completeMultipartUploadRequest.UploadId}
	byteArray, err := util.ToJson(completeMultipartUploadRequest, "parts")

	if err != nil {
		return nil, bce.NewError(err)
	}

	req, err := bce.NewRequest("POST", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName),
		params, bytes.NewReader(byteArray))

	if err != nil {
		return nil, bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")
	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var completeMultipartUploadResponse *CompleteMultipartUploadResponse

	err = json.Unmarshal(res.Body, &completeMultipartUploadResponse)

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

	partETags := make([]PartETag, 0, partCount)

	var waitGroup sync.WaitGroup

	for i := 0; i < partCount; i++ {
		var skipBytes int64 = partSize * int64(i)
		var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

		byteArray := make([]byte, size, size)
		_, err := file.Read(byteArray)

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
			PartData:   byteArray,
		}

		waitGroup.Add(1)

		partETags = append(partETags, PartETag{PartNumber: partNumber})

		go func(partNumber int) {
			defer waitGroup.Done()

			uploadPartResponse, err := c.UploadPart(uploadPartRequest, nil)

			if err != nil {
				panic(err)
			}

			partETags[partNumber-1].ETag = uploadPartResponse.GetETag()
		}(partNumber)
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
			Parts:      partETags,
		}

		completeResponse, err := c.CompleteMultipartUpload(completeMultipartUploadRequest, nil)

		if err != nil {
			panic(err)
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

	req, err := bce.NewRequest("DELETE", c.GetUriPath(objectKey), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return bce.NewError(err)
	}

	_, bceError := c.SendRequest(req, option, false)

	return bceError
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

	req, err := bce.NewRequest("GET", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return nil, bce.NewError(err)
	}

	res, bceError := c.SendRequest(req, option, true)

	if bceError != nil {
		return nil, bceError
	}

	var listMultipartUploadsResponse *ListMultipartUploadsResponse

	err = json.Unmarshal(res.Body, &listMultipartUploadsResponse)

	if err != nil {
		return nil, bce.NewError(err)
	}

	return listMultipartUploadsResponse, nil
}

func (c *Client) setBucketAclFromString(bucketName, acl string, option *bce.SignOption) *bce.Error {
	params := map[string]string{"acl": ""}
	req, err := bce.NewRequest("PUT", c.GetUriPath(""), c.GetBucketEndpoint(bucketName), params, nil)

	if err != nil {
		return bce.NewError(err)
	}

	option = bce.CheckSignOption(option)
	option.AddHeadersToSign("date")

	headers := map[string]string{"x-bce-acl": acl}
	option.AddHeaders(headers)

	_, bceError := c.SendRequest(req, option, true)

	return bceError
}
