package ghttp

import (
        "bytes"
        "io/ioutil"
//        "log"
        "net/http"
        "net/url"
        "net/http/cookiejar"
	"crypto/tls"
        "encoding/json"
)

var gCurCookies []*http.Cookie
var gCurCookieJar *cookiejar.Jar
var gInsecure bool

func gInit() {
        gCurCookies = nil
        //var err error
        gCurCookieJar,_=cookiejar.New(nil)
}

func SetInsecure(val bool) {
	gInsecure=val
}

func Get(urlStr string) (int, string) {
        if gCurCookieJar==nil{
                gInit()
        }
	returnStatus:=0

	gTransport:=&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: gInsecure},}

        httpClient := &http.Client{
                CheckRedirect: nil,
                Jar:           gCurCookieJar,
		Transport: gTransport,
        }

        httpReq, err := http.NewRequest("GET", urlStr, nil)
        httpResp, err := httpClient.Do(httpReq)

        if err != nil {
                //log.Printf("Error=%s\n", err.Error())
		return returnStatus, ""
        }
        //log.Printf("Header=%s", httpResp.Header)
        //log.Printf("Status=%s", httpResp.Status)
	returnStatus=httpResp.StatusCode

        defer httpResp.Body.Close()

        body, errReadAll := ioutil.ReadAll(httpResp.Body)
        if errReadAll != nil {
                //log.Printf("Error=%s\n", errReadAll.Error())
		return returnStatus, ""
        }

        gCurCookies = gCurCookieJar.Cookies(httpReq.URL)

        return returnStatus, string(body)
}

func Do(action string, urlStr string, dataMap map[string]string, jsonSend bool) (int, string) {
        if gCurCookieJar==nil{
                gInit()
        }

	returnStatus:=0

        gTransport:=&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: gInsecure},}

        httpClient := &http.Client{
                CheckRedirect: nil,
                Jar:           gCurCookieJar,
                Transport: gTransport,
        }
	var httpReq *http.Request
	var err error
	if jsonSend==false {
		// URL-encoded
		if dataMap!=nil {
			data:=url.Values{}
			for k,v:=range dataMap {
				data.Add(k, v)
			}
			dataBuffer:=bytes.NewBufferString(data.Encode())
			httpReq, err = http.NewRequest(action, urlStr, dataBuffer)
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}else{
			httpReq, err = http.NewRequest(action, urlStr, nil)
		}
	}else{
		// JSON
		if dataMap!=nil {
			jsonValue,_ := json.Marshal(dataMap)
			dataBuffer:=bytes.NewBuffer(jsonValue)
			httpReq, err = http.NewRequest(action, urlStr, dataBuffer)
			httpReq.Header.Set("Content-Type", "application/json")
		}else{
			httpReq, err = http.NewRequest(action, urlStr, nil)
		}
	}

        httpResp, err := httpClient.Do(httpReq)

        if err != nil {
                //log.Printf("Error=%s\n", err.Error())
		return returnStatus, ""
        }

	returnStatus=httpResp.StatusCode

        defer httpResp.Body.Close()
        body, errReadAll := ioutil.ReadAll(httpResp.Body)
        if errReadAll != nil {
                //log.Printf("Error=%s\n", errReadAll.Error())
		return returnStatus, ""
        }

        gCurCookies = gCurCookieJar.Cookies(httpReq.URL)

        return returnStatus, string(body)
}

