# GOFLY在线客服系统
## 1. 简介

自互联网诞生以来 , 世界上已经有了众多的中小型或者大型网站服务,
当移动化浪潮来临以后,又出现了众多的移动化APP/小程序 . 
为了便利运营者与访客用户进行实时沟通 , 更有效的进行流量转化 , 为开发者或中小站长开发出GOFLY在线客服系统 .
该系统作为访客与网站/APP运营者之间的桥梁 , 提供网页版PC/移动H5的即时通讯聊天的功能 .
该系统极大的帮助了运营者获取详细访客信息 , 实时进行用户画像 ,
获取访客意图 , 扩大运营效果.
该系统同时也极大的提升了访客的用户体验,可以更方便快捷的获取完善有用的信息.

### 1.1 编写目的
本文档为使用说明文档，为产品的使用与维护提供信息基础。
### 1.2	使用对象
本文档的使用对象主要为网站运营者和开发工程师。
### 1.3	产品范围
本系统主要专注于为广大开发工程师和网站运营人员提供网页即时通讯功能,
支持 PC、移动、微信、小程序、APP 等多渠道接入，快速整合自有会员系统 

该系统提供访客端面板pc展示,移动端自适应展示,在访客端实现发送文字消息,发送表情消息,发送图片消息,发送附件文件等功能

该系统提供客服端pc管理系统以及客服端简易H5聊天系统,在客服系统中实现访客到来提醒,访客消息提醒,获取访客来源,获取访客意图,发送客服消息/图片/附件等功能
### 1.4	更新日志
#### v0.6.0

+ 新增访客黑名单功能，可以根据访客id加入黑名单
+ 新增api接口可优雅关闭服务，在守护模式下相当于重启服务
+ 新增客服端搜索客服账号接口
+ 新增微信公众号访客展示是否关注公众号标签
+ 新增客服端首页展示系统公告，超管可以添加管理系统公告
+ 新增系统配置微信模板remark字段，模板消息会展示该字段
+ 修复优化访客表新增real_name字段，客服端首先展示该字段，客服备注姓名存入该字段
* 修复子进程退出次数太多，父守护进程也退出问题
* 修复优化传递商品卡片信息样式效果
* 修复优化访客端图片缩略展示，点击预览大图效果


#### v0.5.9
+ 新增系统配置项，系统管理员权限可以在后台配置本客服系统的标题、关键字、版权等基本信息，以及是否显示注册按钮等
+ 新增商户账号可以查看访客的基本信息，可以查看访客是否绑定了微信公众号
+ 新增商户账号下访客列表展示ip和对应的地址
+ 新增清除访客聊天记录
+ 新增商户配置项，上传微信域名验证文件功能
+ 新增生成微信公众号菜单可视化编辑功能
* 新增微信公众号主动发客服消息接口是否开启配置项
* 修复微信公众号和访客绑定功能，访客关注时新增绑定，访客取关时删除绑定
* 修复访客端超时，另开tab标签等操作时，弹窗确认重新reload页面

#### v0.5.8
+ 新增微信公众号网页oauth授权功能，网页获取微信用户的昵称和头像
+ 新增微信公众号关注时自动回复消息功能，与原来的自动欢迎拆分开
+ 新增访客页面二次跳转到落地域名的配置项
+ 新增微信公众号模板消息，可以给客服、访客发送新消息提醒模板消息
+ 新增微信公众号带参二维码，绑定客服与微信id
+ 新增商户后台配置模板id功能，增加客服消息、访客消息，访客上线三个模板配置
+ 修改客服相关接口的接口前缀
* 修复访客id分割错误问题，访客id连接符修改
* 修复微信公众号access_token获取接口次数超限制问题

#### v0.5.5
+ 新增访客端显示客服在线离线状态，判断客服离线状态
+ 新增自动欢迎内容增加富文本编辑器wangEditor
+ 新增独立链接模式传递访客名称，id和头像信息
+ 新增客服界面展示访客操作系统、浏览器等信息
+ 新增给访客打tag标签功能，可以根据tag标签搜索
+ 新增访客列表搜索功能，根据访客id或者名称进行搜索
+ 修改访客端消息时间的格式
* 修复通知邮件465端口不能发送问题，支持tls的smtp端口
* 修复访客端websocket连接最大次数限制不起作用问题

#### v0.5.1
...


## 2. 产品概述

### 2.1 总体框架
GOFLY客服系统作为WEB即时通讯系统，总体框架如图2-1总体框架图所示。
![Image text](https://gitee.com/taoshihan/go-fly/raw/master/static/images/jiago.jpg)
### 2.2 系统架构
本系统主要使用golang gin框架v1.6.3 , 主要开发语言为golang 1.16,
数据库采用MySQL .基于gin框架实现web服务 ,使用gorm作为数据库关系模型，对外提供http RESTful接口 , 整合jwt-go实现后台权限验证 ,利用websocket实现消息的实时推送.

前端采用vue2.0框架，ui框架使用element-ui，实现了前端的mvvm架构

管理后台提供对在线访客的即时交流，客服人员的管理功能
，对项目的设置功能等

访客端提供对接js接口，可快速嵌入到网页中聊天对话框

### 2.3 模块描述
#### 2.3.1 websocket模块

visitor访客连接ws模块 , 对访客的ws请求进行处理

kefu客服连接ws模块 ,后端客服端连接的websocket模块

websocket公共功能模块，公共功能模块

#### 2.3.1 common公共模块

common模块，提供全局配置公共变量


config模块，获取各个单独配置对象功能

#### 2.3.1 config模块

数据库连接配置文件，可以配置数据库信息

数据库sql文件，数据库的结构和默认数据

ip地址库文件，解析IP为归属地数据文件



#### 2.3.1 controller控制器模块

所有的路由控制器模块，处理具体的路由请求，返回相应的json数据

about模块，处理单页的控制器模块

auth模块，处理权限验证功能

captcha模块，处理验证码功能

chat模块，处理聊天功能

ent模块，处理企业的各项功能

index模块，处理首页部分的功能

ip模块，处理ip黑名单功能

kefu模块，处理客服信息功能

login模块，处理登录请求验证功能

main模块，处理公共请求功能

message模块，处理聊天消息功能

notice模块，处理自动欢迎功能

peer模块，处理webrtc功能

reply模块，处理自动回复功能

reponse模块，公共的结构体定义

role模块，处理角色模块

setting模块，处理设置部分功能

shout模块，处理通知功能模块

visitor模块，处理访客部分的功能

weixin模块，处理对接微信的功能


#### 2.3.1 cmd命令行模块

提供命令行参数中安装系统功能模块,命令行参数中监听端口网站web server主服务模块

root根命令行模块，添加其他cmd命令行模块

install命令行模块，安装系统导入数据库功能

server命令行模块，提供web服务和监听gin框架等服务

stop命令行模块，提供关闭服务功能

#### 2.3.1 middleware中间件模块

IP验证中间件，可以验证当前IP是否在IP黑名单中

jwt权限验证中间件，对后台权限部分进行json web token 验证

logger日志中间件，对系统所有日志的记录和输出

rbac权限验证中间件，对系统后台角色进行权限划分认证

cross跨域中间件，系统对外提供允许跨域请求的功能

language语言中间件，处理多语言功能模块

#### 2.3.1 model模型模块

该模块提供对接数据库获取数据的主要功能，提供数据库关系映射对象，对外提供数据

about模型，单页数据库模型

config模型，配置表数据库模型

ent_config模型，企业各自的数据表配置模型

ipblack模型，ip黑名单表数据库模型

message模型，消息记录数据库模型

models模型，公共连接数据库mysql的功能模型

reply模型，自动回复表数据库模型

role模型，角色表模型

user_role,角色和用户关联关系表数据库模型

user模型，客服表数据库模型

visitor模型，访客表数据库模型

welcome模型，自动欢迎表数据库模型



#### 2.3.1 router路由模块

注册所有的对外路由地址，包括api路由和模板渲染路由

#### 2.3.1 static静态文件模块

前端css/js/images等静态文件，可以使用nginx对外进行提供服务

渲染模板的模板文件，渲染模板时指定的模板文件



#### 2.3.1 tmpl模板模块

框架展示html模板的服务端加载模块，所有的模板路由文件，在此加载静态文件

#### 2.3.1 tools工具模块

file工具类，文件操作工具类

ip工具类，解析ip的工具类

limits工具类，滑动窗口限流工具类

session工具类，模拟实现的session工具类

smtp工具类，发送通知邮件的工具类

sorts工具类，排序功能的工具类

logger工具类，实现logger的日志记录工具类

uuid工具类，生成唯一ID的工具类


## 3. 使用说明

### 3.1 业务流程说明
#### 3.1.1 整体业务流程
详细流程如图3-1 聊天时序图所示

1)访客在打开网页或者打开聊天链接后,发送当前状态到GOFLY客服系统访客登录接口,接口生成唯一ID存入mysql数据库,并返回给前端json信息

2)访客接收到唯一ID后存入本地localStorge,携带唯一ID连接websocket服务,服务端验证唯一ID并建立ws链接,内存存储唯一ID到ws链接的映射关系

3)客服端建立ws连接,内存存储客服ID到ws链接的映射关系

4)访客和客服互发消息,通过客服ID与访客ID关联的ws链接推送给前端

#### 3.1.2 访客业务流程
详细流程如图3-1 聊天时序图所示

通过js引入聊天插件到网页中，访客到来后，会自动把当前访客连接到websokcet，通过该ws连接可以获取到客服端发送来的消息
#### 3.1.3 客服业务流程
详细流程如图3-1 聊天时序图所示
#### 3.1.4 资源访问流程
资源访问采用REST的定义，通过HTTP get/put/post/patch/delete等方法完成。例如：

1）在线访客资源收到REST的get请求，根据请求内容，可以调用相应的服务orm获取资源，也可以将请求转到route进行后续处理获取资源，并返回资源。

2）资源服务层收到REST的put或者patch请求，根据请求内容调用orm或者转到route进行资源更新。
### 3.2 依赖环境说明
#### 3.2.1 数据库MySQL5.5+
1. mysql数据库版本不能低于5.5

2. 安装数据库可以参考mysql官网文档介绍 , 或者使用其他php集成环境中的mysql服务,比如:宝塔面板/phpstudy等

#### 3.2.2 反向代理服务nginx
1. nginx的主要作用是部署域名,提供https/wss加密服务

2. 使用nginx提供支持websocket连接

3. 主要方式为nginx反向代理golang后端监听的端口,并传递指定的http头信息

4. 安装nginx可以参考nginx官网文档介绍 , 或者使用其他php集成环境中的nginx服务,比如:宝塔面板/phpstudy等

### 3.3 服务安装说明

####  3.3.1 下载编译版客服系统
1. 前往官网下载客服zip压缩包

2. 注意区分自己运行系统是32位还是64位,目前压缩包提供的版本为windows 64位/linux 32和64位

3. 配置数据库信息

    在根目录/config/mysql.json文件中配置数据库连接信息,例如:
    
    数据库服务地址:127.0.0.1
    
    数据库端口:3306
    
    数据库民:gofly
    
    数据库用户名:gofly
    
    数据库密码:gofly
    
    配置完数据库信息后执行 ./go-fly install 会清空删除原表数据,导入默认数据库表结构以及基础数据
```php
{
	"Server":"127.0.0.1",
	"Port":"3306",
	"Database":"gofly",
	"Username":"gofly",
	"Password":"gofly"
}
```


    
4. 运行服务

   linux:   ./go-fly server [可选 -p 8082 -d]
   
   windows: go-fly.exe server [可选 -p 8082 -d]
5. 参数说明

   -p 指定端口
   
   -d 是否以daemon守护进程运行
   
   -h 查看帮助
#### 3.3.2  源码安装运行
1. 依赖管理基于go mod,使用非常方便

   go env -w GO111MODULE=on
   
   go env -w GOPROXY=https://goproxy.cn,direct
   
2. 进入gofly源码目录,配置数据库连接信息与上面编译版一致, 导入数据库使用 go run go-fly.go install

3. 源码运行项目, go run go-fly.go server ,参数选项与编译版一致

4. 源码打包为二进制文件, go build go-fly.go 会生成go-fly可以执行文件,执行就可以直接运行编译版.


#### 3.3.3 nginx部署配置域名

访问：

1.参考支持https的部署示例ssl_系列的指令部分 , 注意反向代理的端口号和证书地址 , 不使用https也可以访问 , 只是不会有浏览器通知弹窗

2.尽量按照下面的配置处理, 配置独立域名或者二级域名, 不建议在主域名加端口访问, 不建议主域名加目录访问 

3.参考注意支持websocket部分的配置 ,下面两项是支持ws的配置

proxy_set_header Upgrade $http_upgrade;

proxy_set_header Connection "upgrade";

4.透传客户端IP,默认反向代理后，后端获取不到客户端IP，通过这个http头可以透传客户端IP给后端服务

proxy_set_header X-Real-IP $remote_addr;

```php
server {
       listen 443 ssl http2;
        ssl on;
        ssl_certificate   conf.d/cert/4263285_gofly.sopans.com.pem;
        ssl_certificate_key  conf.d/cert/4263285_gofly.sopans.com.key;
        ssl_session_timeout 5m;
        ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_prefer_server_ciphers on;
        #listen          80; 
        server_name  gofly.sopans.com;
        access_log  /var/log/nginx/gofly.sopans.com.access.log  main;
        location / {
                proxy_pass http://127.0.0.1:8081;
                proxy_http_version 1.1;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "upgrade";
        }
}
server{
       listen 80;
        server_name  gofly.sopans.com;
        access_log  /var/log/nginx/gofly.sopans.com.access.log  main;
        location / {
                proxy_pass http://127.0.0.1:8081;
                proxy_http_version 1.1;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection "upgrade";
        }
}
```
### 3.4 访客接入

#### 3.4.1 超链接模式

超链接模式是最简单易懂的方式 , 通过下面的路径可以查询到本商户的配置信息

访问路径:后台==>设置==>商户服务==>网页部署

聊天面板路径:

[域名]/chatIndex?kefu_id=[商户名称]&ent_id=[商户ID]

#### 3.4.2 悬浮窗模式

悬浮窗模式是符合当前主流客服系统的接入方式 , 通过引入一个js文件 , 可自动弹出咨询客服按钮，访客点击按钮即可在线聊天。

```php
<!--对接客服代码-->
<script src="[域名]/static/js/gofly-front.js"></script>
<script>
    GOFLY.init({
        GOFLY_URL:"[域名]",//域名,必填
        GOFLY_KEFU_ID: "kefu2",//商户名称,必填
        GOFLY_ENT: "2",//商户ID,必填
        GOFLY_BTN_TEXT: "GOFLY-PRO LIVE CHAT",//按钮文件,选填
        GOFLY_EXTRA:{
            "visitorProduct":{
                "title":"GOFLY客服系统商务版",
                "price":"1000元",
                "img":"https://jd.sopans.com/static/upload/2021April/9efd50ba5d97f1136137ed252576a95e.png",
                "url":"",
            }
        }
    })
</script>
<!--//对接客服代码-->
```

商品链接格式

product[{"title":"GOFLY客服系统商务版","price":"100000元","img":"https://jd.sopans.com/static/upload/2021April/9efd50ba5d97f1136137ed252576a95e.png","url":""}]
### 版权声明

gofly客服系统为GOFLY独自开发,GOFLY拥有系统的全部版权,其他人使用请务必联系本人,获取授权后方可使用.