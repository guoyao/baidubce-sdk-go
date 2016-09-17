package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/guoyao/baidubce-sdk-go/bce"
)

var credentials = bce.NewCredentials(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK"))

//var bceConfig = bce.NewConfig(credentials)
var bceConfig = &bce.Config{
	Credentials: credentials,
	Checksum:    true,
	//Protocol:    "https",
}

var bceClient = bce.NewClient(bceConfig)

func init() {
	bceClient.SetDebug(true)

	/*
		bceConfig.Endpoint = "baidubce-sdk-go.bj.bcebos.com"
		bceConfig.ProxyHost = "agent.baidu.com"
		bceConfig.ProxyPort = 8118
		bceConfig.MaxConnections = 6
		bceConfig.Timeout = 6 * time.Second
	*/
}

func getSessionToken() {
	req := bce.SessionTokenRequest{
		DurationSeconds: 600,
		Id:              "ef5a4b19192f4931adcf0e12f82795e2",
		AccessControlList: []bce.AccessControlListItem{
			bce.AccessControlListItem{
				Service:    "bce:bos",
				Region:     "bj",
				Effect:     "Allow",
				Resource:   []string{"baidubce-sdk-go/*"},
				Permission: []string{"READ"},
			},
		},
	}

	sessionTokenResponse, err := bceClient.GetSessionToken(req, nil)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(sessionTokenResponse)
	}
}

func RunBceExamples() {
	getSessionToken()
}
