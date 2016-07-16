package util

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

const URI string = "http://bos.cn-n1.baidubce.com/v1/example/测试"

func TestGetURIPath(t *testing.T) {
	expected := "/v1/example/测试"
	path := GetURIPath(URI)

	if path != expected {
		t.Error(FormatTest("GetURIPath", path, expected))
	}
}

func TestURIEncodeExceptSlash(t *testing.T) {
	expected := "/v1/example/%E6%B5%8B%E8%AF%95"
	path := GetURIPath(URI)
	path = URIEncodeExceptSlash(path)

	if path != expected {
		t.Error(FormatTest("URIEncodeExceptSlash", path, expected))
	}
}

func TestGetMD5(t *testing.T) {
	expected := "de22e061b93b832dd8af907ca9002fd7"
	result := GetMD5("baidubce-sdk-go", false)

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}
}

func TestTimeToUTCString(t *testing.T) {
	expected := "2015-11-16T07:33:15Z"
	datetime, _ := time.Parse(time.RFC1123, "Mon, 16 Nov 2015 15:33:15 CST")
	utc := TimeToUTCString(datetime)

	if utc != expected {
		t.Error(FormatTest("TimeToUTCString", utc, expected))
	}
}

func TestTimeStringToRFC1123(t *testing.T) {
	expected := "Mon, 16 Nov 2015 07:33:15 UTC"
	result := TimeStringToRFC1123("2015-11-16T07:33:15Z")

	if result != expected {
		t.Error(FormatTest("TimeStringToRFC1123", result, expected))
	}
}

func TestHostToURL(t *testing.T) {
	expected := "http://bj.bcebos.com"
	host := "bj.bcebos.com"
	url := HostToURL(host, "http")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}

	host = "http://bj.bcebos.com"
	url = HostToURL(host, "http")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}
}

func TestToCanonicalQueryString(t *testing.T) {
	const expected = "text10=test&text1=%E6%B5%8B%E8%AF%95&text="
	params := map[string]string{
		"text":   "",
		"text1":  "测试",
		"text10": "test",
	}
	encodedQueryString := ToCanonicalQueryString(params)

	if encodedQueryString != expected {
		t.Error(FormatTest("ToCanonicalQueryString", encodedQueryString, expected))
	}
}

func TestToCanonicalHeaderString(t *testing.T) {
	expected := strings.Join([]string{
		"content-length:8",
		"content-md5:0a52730597fb4ffa01fc117d9e71e3a9",
		"content-type:text%2Fplain",
		"host:bj.bcebos.com",
		"x-bce-date:2015-04-27T08%3A23%3A49Z",
	}, "\n")

	canonicalHeader := ToCanonicalHeaderString(getHeaders())

	if canonicalHeader != expected {
		t.Error(FormatTest("ToCanonicalHeaderString", canonicalHeader, expected))
	}
}

func TestURLEncode(t *testing.T) {
	expected := "test-%E6%B5%8B%E8%AF%95"
	result := URLEncode("test-测试")

	if result != expected {
		t.Error(FormatTest("URLEncode", result, expected))
	}
}

func TestSliceToLower(t *testing.T) {
	expected := "name age"
	arr := []string{"Name", "Age"}
	SliceToLower(arr)

	result := fmt.Sprintf("%s %s", arr[0], arr[1])

	if result != expected {
		t.Error(FormatTest("SliceToLower", result, expected))
	}
}

func TestMapKeyToLower(t *testing.T) {
	expected := "name gender"
	m := map[string]string{"Name": "guoyao", "Gender": "male"}
	MapKeyToLower(m)

	result := ""

	if _, ok := m["name"]; ok {
		result += "name"
	}

	if _, ok := m["gender"]; ok {
		result += " gender"
	}

	if result != expected {
		t.Error(FormatTest("MapKeyToLower", result, expected))
	}
}

func TestToMap(t *testing.T) {
	p := struct {
		Name   string
		Age    int
		Gender string
	}{"guoyao", 10, "male"}

	m, err := ToMap(p, "Name", "Age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		expected := "guoyao:10"
		result := fmt.Sprintf("%s:%v", m["Name"], m["Age"])

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}
}

func TestToJson(t *testing.T) {
	p := struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	}{"guoyao", 10, "male"}

	byteArray, err := ToJson(p, "name", "age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		expected := "{\"age\":10,\"name\":\"guoyao\"}"
		result := string(byteArray)

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}
}

func TestCheckFileExists(t *testing.T) {
	expected := true
	result := CheckFileExists("util_test.go")

	if result != expected {
		t.Error(FormatTest("CheckFileExists", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func getHeaders() map[string]string {
	header := map[string]string{
		"Host":           "bj.bcebos.com",
		"Content-Type":   "text/plain",
		"Content-Length": "8",
		"Content-Md5":    "0a52730597fb4ffa01fc117d9e71e3a9",
		"x-bce-date":     "2015-04-27T08:23:49Z",
	}

	return header
}
