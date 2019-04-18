package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	//"strconv"
)

type Body struct {
	ShopId string `json:"shopId"`
	DevSN  string `json:"devSN"`
	ApSN   string `json:"apSN"`
}

type respData struct {
	Code int  `json:"code"`
	Data Data `json:"data"`
}

type Data struct {
	ApName     string `json:"apName"`
	FailReason string `json:"failReason"`
	RunTime    int    `json:"runTime"`
}

const debug bool = false
const warning string ="×××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××"

func main() {
	body := Body{"768313", "210235A1JTB15C000019", "219801A0WA9163Q09539"}
	for i := 1; ; i++ {
		//baidu()

		result, err := getFailReason(body)
		if err != nil {
			log.Printf("Get reason fail,%v!\n", err)
		} else {
			data1 := new(respData)
			json.Unmarshal(result, data1)
			if debug {
				log.Printf("data is %v,result is %v\n", data1, string(result))
			}
			if data1.Data.FailReason == " Kernel exception reboot\r\n" {
				log.Println("AP fail for kernel exception reboot!")
				temp := data1.Data.RunTime
				if temp<300{
					log.Printf("\n%s\n%s\n%s\n",warning,warning,warning)
					panic(2)
				}
				log.Printf("AP run for %v,as %vd %vh %vm %vs!", temp, (temp)/(24*60*60), ((temp)%(24*60*60))/3600, (((temp)%(24*60*60))%3600)/60, ((((temp) % (24 * 60 * 60)) % 3600) % 60))
			} else {
				log.Printf("AP fail for %v! Get count %v.\n", data1.Data.FailReason, i)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func getFailReason(body Body) ([]byte, error) {
	bodyJson, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://oasis.h3c.com/v3/apmonitor/version/1.0/maintain/getApBasicInfoDetail", bytes.NewBuffer(bodyJson))
	req.Header = map[string][]string{
		"accept":       {"application/json, text/plain, */*"},
		"User-Agent":   {"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36"},
		"content-type": {"application/json"},
		"Host":         {"oasis.h3c.com"},
		"cookie":       {"lang=cn; connect.sid=s%3AGjaMZmbpo4iWv6FCxgRD8E-PSEkWmXA5.ry4tA8hkmGT6ZqcPtc2WVj6pM1HRGmwZXNH%2BlVQ4Gyk"},
	}
	fixedURL, err := url.Parse("http://y13709:yy.1220.F@devproxy.h3c.com:8080")
	transport := &http.Transport{
		Proxy: http.ProxyURL(fixedURL),
	}
	client := &http.Client{Transport: transport}
	if debug {
		temp2 := make([]byte, 1024)
		temp2, _ = httputil.DumpRequest(req, true)
		log.Println(string(temp2))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if debug {
		temp := make([]byte, 1024)
		temp, _ = httputil.DumpResponse(resp, true)
		log.Println(string(temp))
	}
	body2, _ := ioutil.ReadAll(resp.Body)
	return body2, err
}

func baidu() {
	req, _ := http.NewRequest("GET", "https://www.google.com", nil)
	fixedURL, err := url.Parse("http://y13709:yy.1220.F@devproxy.h3c.com:8080")
	transport := &http.Transport{
		Proxy: http.ProxyURL(fixedURL),
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err, "baidu")
		return
	}
	defer resp.Body.Close()
	httputil.DumpResponse(resp, true)
	fmt.Printf("%v", resp)
}
