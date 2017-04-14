## 抽象：
想象有个人在浏览器中点击操作。
也就是有request对象，和response对象。</br>

要做的事==模拟一个客户端

看了《HTTP权威指南》,发现没啥卵用。除了重新复习下头部信息,cookie操作。

## 代码实现
刚开始还是蛮模糊的,隐约知道怎么做。但是代码怎么写，还没有想好。</br>
然后就上谷歌，看有没有人做过。按照前人的指导，我就慢慢做出来了。</br>
模拟成功的判断的依据是，post后的页面，err_no=0</br>
刚开始，其实我也没啥把握。尤其是需要post那么多的数据。去哪找。后来，想无论如何，还是完成吧，好歹知道到哪一步不成功，有个交代。</br>

## 步骤分解
网上都有教程，讲地不是有的比较罗嗦，没重点。要么就是太简单。经过反复查看他们的文章。我简单地写下:</br>
讲许多，`重要的是作品完成之前的思考`。</br>
许多人都没讲到这个。所以大家可以去我下面给的链接，里面有讲</br>
怎么去着手模拟百度登录。</br>

可以用IE9,chrome,firefox;</br>
我觉得firefox比较适合，顺手。但是也不否认，chrome的价值。</br>
我是开着两个浏览器，F12看http的网络请求。</br>

大家都知道,网上的东西不一定可信，并且斟酌信息的正确也比较难。找东西也不是简单

我在这写，遇到的几个难点，其他的都有前辈给解决掉了。

### 1.RSA加密
之前没有做过，拿到手上一时半会没反应。继续谷歌。

### 2.百度的中文验证码
查了golang的库，有识别字母数字的，但是没有识别中文的。 有用cgo,但是我也不熟悉c，c++.再着主要目的是成功模拟登录。我就肉眼识别了。

### 3.用go调用js
这个之前也没有做过，相当于刚上手。</br>
有库:`github.com/robertkrimen/otto`

###4. 最重要的是理清流程。
分析数据来源，有哪几个步骤。

## 步骤

### 1.获取cookie:BAIDUID

### 2.获取token
需要本地生成gid,调用js `src/utils/pswd.js`</br>
参考代码:

```js
 guidRandom()
 ```
### 3.获取public key
登录加密的

参考代码:
`src/utils/pswd.js`
```js
PasswordEncrypt(pswd, pubkey)
 ```
### 4.获取CodeString
获取验证码。我看七夜前辈，有确认验证码是否正确的步骤，我也写了。不知道，有没有用

### 5.post表单
另外，post数据，有固定的，也有变化的，我把变化的拿出来
```golang
postDict["codestring"] = this.CodeString//验证码的url的参数
postDict["verifycode"] = this.ReadPerson//验证码
postDict["token"] = this.Token//get到的
postDict["rsakey"] = this.PublicKeyResponse.Key //RSA的key
postDict["gid"] = this.Gid//生成的
postDict["password"] = this.RsaPasswd //public KEY加密的
```

-------
以上就是模拟百度登录的。
获取数据的url的参数，都是组合的。就是说，有的数据是要`找`的。我就不写了，不方便写，不是固定的。

最后，流程都有了，就是写代码的事情了。要参考golang的库，鼓起勇气看函数的提供，看函数的讲解。</br>
七拼八凑，好歹是写完了。

### 模拟百度登录
弄了一周多点时间的百度登录

参考文档：
- http://www.crifan.com/use_ie9_f12_to_analysis_the_internal_logical_process_of_login_baidu_main_page_website/
- http://www.cnblogs.com/qiyeboy/p/5722424.html 
- http://bbs.125.la/forum.php?mod=viewthread&tid=13881883&page=1#pid9286935   
