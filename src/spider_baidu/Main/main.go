package main

import (
	"fmt"
	"os"
	spider "spider_baidu"
	"utils"
)

const (
	Login_url string = "https://passport.baidu.com/v2/api/?login"
	UID       string = "http://www.baidu.com"
	IMG       string = ""
)

var spiderYello *spider.Spider = spider.NewSpider()

func main() {
	//第一步获取cookie:BAIDUID
	fmt.Println("<<获取cookie:BAIDUID>>")
	spiderYello.Get_client(UID)
	//第二步获取token
	//设置GID
	fmt.Println("<<获取token>>")
	setGID()
	TOKEN := spiderYello.MakeTokenUrl()
	tokenBody := spiderYello.Get_client(TOKEN)
	spiderYello.GetToken(tokenBody)
	//获取public key
	fmt.Println("<<获取public key>>")
	makePublicKey()
	//设置post参数
	setRsaPasswd()
	fmt.Println("<<获取CodeString>>")
	setPic()
	//检查输入的验证码
	checkImgUrl := spiderYello.CheckImg()
	chkeckImgBody := spiderYello.Get_client(checkImgUrl)
	fmt.Println("验证码返回的:", chkeckImgBody)
	//spiderYello.GetToken(tokenBody)
	//第三步登录BAIDU
	dic := spiderYello.MakePost()
	postBody := spiderYello.Post_client(Login_url, dic)
	fmt.Println("post后返回的页码是:", postBody)
}

func makePublicKey() {
	getPublicUrl := spiderYello.MakePublicUrl()
	pub := spiderYello.Get_client(getPublicUrl)
	spiderYello.GetPublicKey(pub)
}

func setRsaPasswd() {
	//获取gid,获取rsa密码
	fi, err := os.Open("../../utils/pswd.js")
	if err != nil {
		panic(err)
	}
	spiderYello.RsaPasswd = utils.MakePasswd(fi, spiderYello.PublicKeyResponse.Pubkey, "qweASDzxc!@#")
	defer fi.Close()
}

func setGID() {
	//获取gid,获取rsa密码
	fi, err := os.Open("../../utils/pswd.js")
	if err != nil {
		panic(err)
	}
	spiderYello.Gid = utils.MakeGid(fi)
	defer fi.Close()
	fmt.Printf("Gid %s\n", spiderYello.Gid)
}

func setPic() {
	//获取CodeString
	codeStringURL := spider.MakeCodeStringUrl(spiderYello.Token)
	codeBody := spiderYello.Get_client(codeStringURL)
	spiderYello.MakeCodeString(codeBody)
	fmt.Println("ImgUrlCode", spiderYello.CodeString)
	//抓取图片
	spider.MakeImg(spiderYello.CodeString)
	//等待识别图片
	spiderYello.ReadPic()
}
