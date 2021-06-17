package main

import (
	"coremail_address_list_export/commonutils"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	help       bool
	httpClient *http.Client
	outpath    string = "email.txt"
	outType    string = "xlsx"
	timeout    int    = 1
	proxy      string = "http://127.0.0.1:8080"
	sid        string
	cookie     string
	headers    map[string]string = map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
		"Content-Type": "text/x-json; charset=UTF-8",
		"Cookie":       fmt.Sprintf("Coremail=%s", cookie),
	}
	mailDomain      string
	mailGroupUrl    string                 = fmt.Sprintf("%s/coremail/s/json?sid=%s&func=oab%%3AgetDirectories", mailDomain, sid)
	mailDataUrl     string                 = fmt.Sprintf("%s/coremail/s/json?sid=%s&func=oab%%3AlistEx", mailDomain, sid)
	columns         []string               = []string{"Name", "Email", "Mobile"}
	columnLength    map[string]int         = map[string]int{"Name": 30, "Email": 20, "Mobile": 12}
	headStyle       *commonutils.XlsxStyle = commonutils.NewXlsxStyle("center", "00A2A5A1", "Verdana", 13)
	cellCenterStyle *commonutils.XlsxStyle = commonutils.NewXlsxStyle("center", "FFFFFFFF", "宋体", 11)
	cellLeftStyle   *commonutils.XlsxStyle = commonutils.NewXlsxStyle("left", "FFFFFFFF", "宋体", 11)
	column_tags     []commonutils.XlsxCol  = []commonutils.XlsxCol{commonutils.XlsxCol{0, "Index", 5.5, headStyle.Style, cellCenterStyle.Style}, commonutils.XlsxCol{1, "Name", 15.0, headStyle.Style, cellLeftStyle.Style}, commonutils.XlsxCol{2, "Email", 24.5, headStyle.Style, cellLeftStyle.Style}, commonutils.XlsxCol{3, "Mobile", 10.0, headStyle.Style, cellLeftStyle.Style}}
	result_list     []interface{}
)

type MailGroup struct {
	Dn      string
	AttrIds []string
}

type MailReqData struct {
	Dn                       string
	ReturnAttrs              []string
	Start                    int
	Limit                    int
	DefaultReturnMeetingRoom bool
}

type MailListData struct {
	True_name     string
	Mobile_number string
	Duty          string
	Gender        string
	Email         string
}

type MailList struct {
	Code  string
	Var   []MailListData
	Total int
}

type MailGroupRespVar struct {
	Id   string
	Name string
	Ou   []MailGroupRespVar
}

type MailGroupRespData struct {
	Code string
	Var  MailGroupRespVar
}

func FetchMailDatas(id string) {
	httpClient, _ = commonutils.NewHttpClient(timeout, proxy)
	dn := "a"
	if id != "" {
		dn = fmt.Sprintf("a/%s", id)
	}
	data := &MailReqData{
		Dn:                       dn,
		ReturnAttrs:              []string{"true_name", "email", "mobile_number", "duty", "gender"},
		Start:                    0,
		Limit:                    1000000,
		DefaultReturnMeetingRoom: false,
	}
	json_data, _ := json.Marshal(data)
	post_data := strings.ToLower(string(json_data))
	post_data = strings.Replace(post_data, "returnattrs", "returnAttrs", 1)
	resp, err := commonutils.HttpPost(mailDataUrl, post_data, headers, timeout, httpClient)
	if err == nil {
		respDatas := &MailList{}
		err3 := json.Unmarshal([]byte(resp["respBody"]), &respDatas)
		if err3 == nil {
			fmt.Printf("Total: %d\n", respDatas.Total)
			if respDatas.Var != nil && len(respDatas.Var) > 0 {
				for _, v := range respDatas.Var {
					fmt.Printf("%s,%s,%s\n", v.True_name, v.Email, v.Mobile_number)
					res := make(map[string]string, 3)
					res["Name"] = v.True_name
					res["Email"] = v.Email
					res["Mobile"] = v.Mobile_number
					// fmt.Println(res["Url"], res["StatusCode"], res["Title"])
					result_list = append(result_list, res)
				}
			}
		}
	}
}

func FetchMailGroups(respVarDatas MailGroupRespVar) {
	if respVarDatas.Ou != nil && len(respVarDatas.Ou) > 0 {
		for _, v := range respVarDatas.Ou {
			fmt.Printf("\nGroup Id: %s,Group Name: %s\n", v.Id, v.Name)
			FetchMailDatas(v.Id)
			if v.Ou != nil && len(v.Ou) > 0 {
				FetchMailGroups(v)
			}
		}
	}
}

func main() {
	flag.StringVar(&mailDomain, "u", "", "http://mail.xxx.com")
	flag.IntVar(&timeout, "t", 7, "http request timeout")
	flag.StringVar(&proxy, "proxy", "", "proxy: http://ip:port, socks5://ip:port")
	flag.StringVar(&sid, "sid", "", "Cookie.Coremail.sid")
	flag.StringVar(&cookie, "cookie", "", "Cookie.Coremail")
	flag.StringVar(&outType, "oT", "xlsx", "out file type:txt,csv,xlsx")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}
	if mailDomain == "" || sid == "" || cookie == "" {
		fmt.Println("param can't empty!,-u -sid -cookie")
		flag.Usage()
		os.Exit(0)
	} else {
		headers = map[string]string{
			"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
			"Content-Type": "text/x-json; charset=UTF-8",
			"Cookie":       fmt.Sprintf("Coremail=%s", cookie),
		}
		mailGroupUrl = fmt.Sprintf("%s/coremail/s/json?sid=%s&func=oab%%3AgetDirectories", mailDomain, sid)
		mailDataUrl = fmt.Sprintf("%s/coremail/s/json?sid=%s&func=oab%%3AlistEx", mailDomain, sid)
	}
	fmt.Printf("\nmailDomain: %s\n", mailDomain)
	fmt.Printf("sid: %s\n", sid)
	fmt.Printf("cookie: %s\n\n", cookie)
	httpClient, _ = commonutils.NewHttpClient(timeout, proxy)
	data := &MailGroup{Dn: "a", AttrIds: []string{"email"}}
	json_data, err1 := json.Marshal(data)
	if err1 != nil {
		fmt.Println(err1)
	}
	resp, err := commonutils.HttpPost(mailGroupUrl, strings.ToLower(string(json_data)), headers, timeout, httpClient)
	if err == nil {
		respDatas := &MailGroupRespData{}
		err3 := json.Unmarshal([]byte(resp["respBody"]), &respDatas)
		if err3 == nil {
			fmt.Println(respDatas.Var.Name)
			FetchMailDatas("")
			if respDatas.Var.Ou != nil && len(respDatas.Var.Ou) > 0 {
				FetchMailGroups(respDatas.Var)
			}
		}
		fmt.Printf("\nDone!\nfetch email count: %v\n", len(result_list))
		if len(result_list) > 0 {
			i := 1
			for _, v := range result_list {
				v.(map[string]string)["Index"] = strconv.Itoa(i)
				i++
			}
			outpath = time.Now().Format("Email_2006-01-02_150405") + "." + outType
			isSaved := false
			if outType == "xlsx" {
				isSaved = commonutils.WriteFileLinesToExcel(column_tags, outpath, result_list)
			} else if outType == "csv" {
				isSaved = commonutils.WriteFileLinesBySplitChar(columns, outpath, result_list, ",")
			} else {
				isSaved = commonutils.WriteFileLinesBySplitChar(columns, outpath, result_list, " ")
			}
			if isSaved {
				fmt.Printf("save to file %s ……\n", outpath)
			} else {
				fmt.Println("save result failed!")
			}
		}
	}
}
