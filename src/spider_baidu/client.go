package spider_baidu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

/*第一步获取BAIDU_UID  get请求
第二步获取TOKEN            get请求
第三步登录BAIDU             post请求
*/

var httpReq *http.Request

type Spider struct {
	http_client http.Client
	gCurCookies []*http.Cookie
	Token       string
	Gid         string
	RsaPasswd   string
	PublicKeyResponse
	ReadPerson string
	CodeString string
}

type PublicKeyResponse struct {
	Errno  string `json:"errno"`
	Msg    string `json:"msg"`
	Pubkey string `json:"pubkey"`
	Key    string `json:"key"`
}

func NewSpider() *Spider {
	this := &Spider{}
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("设置cookieJar出错", err.Error())
		return nil
	}
	this.http_client.Jar = cookieJar
	return this
}

func (this *Spider) MakeTokenUrl() string {
	var tokenUrl string
	time := time.Now().UnixNano()
	times := strconv.FormatInt(time, 10)
	tokenUrl = "https://passport.baidu.com/v2/api/?getapi&tpl=mn&apiver=v3&tt=" + times[:len(times)-6] + "&class=login&gid=" + this.Gid + "&logintype=dialogLogin&callback=bd__cbs__bmlhf3"
	fmt.Println("获取tokenURL", tokenUrl)
	return tokenUrl
}

//第一步请求Cookie
func (this *Spider) Get_client(url string) string {
	//新建一个request
	httpReq, err := http.NewRequest("GET", url, nil)
	//defer httpReq.Body.Close()
	if err != nil {
		fmt.Println("新建一个request失败", err.Error())
		return ""
	}
	//客户端执行request
	httpResp, err := this.http_client.Do(httpReq)
	defer httpResp.Body.Close()
	if err != nil {
		fmt.Println("get百度网页失败", err.Error())
		return ""
	}
	//获取网页内容
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("读取网页内容失败", err.Error())
		return ""
	}
	fmt.Println("需要发送的cookie")
	for k, v := range this.http_client.Jar.Cookies(httpReq.URL) {
		fmt.Printf("%d   %+v\n", k+1, v)
	}
	fmt.Println("==========")
	//Response回复的cookie
	responseCookie := httpResp.Cookies()

	fmt.Println("回复的cookie")
	for k, v := range responseCookie {
		this.gCurCookies = append(this.gCurCookies, v)
		fmt.Printf("%d   %+v\n", k+1, v)
	}
	fmt.Println("<<==========")
	//设置cookie管理器
	this.http_client.Jar.SetCookies(httpReq.URL, responseCookie)
	for _, v := range responseCookie {
		this.gCurCookies = append(this.gCurCookies, v)
	}
	this.Console_cookie()
	return string(body)
}

func (this *Spider) Post_client(strUrl string, postDict map[string]string) string {
	postValues := url.Values{}
	for postKey, PostValue := range postDict {
		postValues.Set(postKey, PostValue)
	}
	postDataStr := postValues.Encode()
	postDataBytes := []byte(postDataStr)
	postBytesReader := bytes.NewReader(postDataBytes)
	httpReq, _ = http.NewRequest("POST", strUrl, postBytesReader)
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//客户端执行request
	httpResp, err := this.http_client.Do(httpReq)
	if err != nil {
		fmt.Println("get百度网页失败", err.Error())
		return ""
	}

	defer httpResp.Body.Close()
	body, errReadAll := ioutil.ReadAll(httpResp.Body)
	if errReadAll != nil {
		panic(errReadAll)
	}
	fmt.Println("Post需要发送的cookie")
	for k, v := range this.http_client.Jar.Cookies(httpReq.URL) {
		fmt.Printf("%d   %+v\n", k+1, v)
	}
	fmt.Println("===============")
	//Response回复的cookie
	responseCookie := httpResp.Cookies()
	//设置cookie管理器
	this.http_client.Jar.SetCookies(httpReq.URL, responseCookie)

	fmt.Println("Post回复的cookie")
	for k, v := range responseCookie {
		this.gCurCookies = append(this.gCurCookies, v)
		fmt.Printf("%d   %+v\n", k+1, v)
	}
	fmt.Println("<<==========")
	this.Console_cookie()
	return string(body)
}

func (this *Spider) ReadPic() {
	fmt.Scanln(&this.ReadPerson)
}

func (this *Spider) MakePublicUrl() string {
	var getPublicUrl string
	time := time.Now().UnixNano()
	times := strconv.FormatInt(time, 10)
	getPublicUrl = "https://passport.baidu.com/v2/getpublickey" + "?token=" + this.Token + "&tpl=mn&apiver=v3&tt=" + times[:len(times)-6] + "&gid=" + this.Gid + "&callback=bd__cbs__d3en7h"
	fmt.Println("PublicURL:", getPublicUrl)
	return getPublicUrl
}

func MakeCodeStringUrl(token string) string {
	var codeStringUrl string
	time := time.Now().UnixNano()
	times := strconv.FormatInt(time, 10)
	codeStringUrl = `https://passport.baidu.com/v2/api/?logincheck&` + "?token=" + token + `&tpl=mn&apiver=v3&tt=` + times[:len(times)-6] + `&sub_source=leadsetpwd&username=15868478830&isphone=false&dv=&callback=bd__cbs__6fq7ta`
	fmt.Println("CodeStringUrl :", codeStringUrl)
	return codeStringUrl
}

func (this *Spider) MakeCodeString(codeBody string) {
	cstring := strings.Split(codeBody, `"codeString" : "`)
	cstring2 := strings.Split(cstring[1], `",        "vcodetype"`)
	this.CodeString = cstring2[0]
}

func MakeImg(codeString string) {
	var imgUrl string
	imgUrl = "https://passport.baidu.com/cgi-bin/genimage?" + codeString
	//create file
	out, err := os.Create("erwei.png")
	defer out.Close()
	//get html
	resp, err := http.Get(imgUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//read html's body
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		panic(err)
	}
}

//打印cookie
func (this *Spider) Console_cookie() {
	fmt.Println("收藏的cookie")
	for k, v := range this.gCurCookies {
		fmt.Printf(" %d cookie %+v\n", k+1, v)
	}
	fmt.Println("<<====================")
}

func (this *Spider) CheckImg() string {
	var checkImgUrl string
	time := time.Now().UnixNano()
	times := strconv.FormatInt(time, 10)
	checkImgUrl = `https://passport.baidu.com/v2/?checkvcode&token=` + this.Token + `&tpl=mn&apiver=v3&tt=` + times[:len(times)-6] + `&verifycode=` + this.ReadPerson + `&codestring=` + this.CodeString + `&callback=bd__cbs__fkdrws`
	fmt.Println("checkImgUrl", checkImgUrl)
	return checkImgUrl
}

func (this *Spider) GetToken(body string) {
	cstring := strings.Split(body, `"token" : "`)
	cstring2 := strings.Split(cstring[1], `",        "cookie"`)
	this.Token = cstring2[0]
	fmt.Println(" token", this.Token)
}

func (this *Spider) GetPublicKey(pub string) {
	fmt.Println("pub", pub)
	pub1 := strings.Split(pub, "bd__cbs__d3en7h(")
	pub2 := strings.Split(pub1[1], ")")
	dd := strings.Replace(pub2[0], `\n`, "", -1)
	//dd = strings.Replace(dd, `\/`, "", -1)
	dd = strings.Replace(dd, `'`, `"`, -1)

	var p PublicKeyResponse
	if err := json.Unmarshal([]byte(dd), &p); err != nil {
		fmt.Println("marshal err1111:", err.Error())
		return
	}
	fmt.Printf("Marshal %#v\n", p)
	this.PublicKeyResponse = p
}

func (this *Spider) MakePost() map[string]string {
	postDict := map[string]string{}
	time := time.Now().UnixNano()
	times := strconv.FormatInt(time, 10)
	postDict["staticpage"] = "https://www.baidu.com/cache/user/html/v3Jump.html"
	postDict["charset"] = "utf-8"
	postDict["tpl"] = "mn"
	postDict["subpro"] = ""
	postDict["apiver"] = "v3"
	postDict["tt"] = times[:len(times)-6]
	postDict["safeflg"] = "0"
	postDict["u"] = "https://www.baidu.com/"
	postDict["isPhone"] = "false"
	postDict["quick_user"] = "0"

	postDict["detect"] = "1"
	postDict["logintype"] = "dialogLogin"
	postDict["logLoginType"] = "pc_loginDialog"
	postDict["idc"] = ""
	postDict["loginmerge"] = "true"
	postDict["mem_pass"] = "on "
	postDict["splogin"] = "rate"
	postDict["countrycode"] = ""
	postDict["crypttype"] = "12"

	postDict["dv"] = ""
	postDict["ppui_logintime"] = "59702"
	postDict["callback"] = "parent.bd__pcbs__nbe28l"
	postDict["username"] = "15868478830"
	postDict["codestring"] = this.CodeString
	postDict["verifycode"] = this.ReadPerson
	postDict["token"] = this.Token
	postDict["rsakey"] = this.PublicKeyResponse.Key
	postDict["gid"] = this.Gid
	postDict["password"] = this.RsaPasswd
	fmt.Printf("Post内容:%+v \n", postDict)
	return postDict
}
