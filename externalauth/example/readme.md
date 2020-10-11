# External Auth Demo代码

这份代码展示了怎么使用 github.com/herb-go/externalauth 库进行第三方登陆工作。

## 使用方式

1. 将config.json.exmaole复制为config.json,并根据响应的第三方登录服务商提供的信息填写响应的配置
2. 配置域名/host文件，使得可以通过域名打开测试站。注意部分登录方式需要使用HTTPS证书
3. 将域名信息更新到Config.json文件的Host字段
在各个第三方登录服务商处登记回调地址，格式为 "主域名/auth/auth/服务id"。服务id为json中各配置的主键
4. 在浏览器中打开测试站，点击各个链接进行登录测试

## 支持的服务一览

* Github:github 帐号登录
* WechatworkQR:微信企业号扫码登录
* WechatworkOauth:微信企业号微信端登录
* Windowslive:微软windows live 帐号登录
* Wechatmp:微信公众号微信端登录
* Facebook:facebook 帐号登录