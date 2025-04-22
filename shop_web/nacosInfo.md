## user-web.json
Group：dev
```
{
  "name": "user-web",
  "port": 8021,
  "host":"本机IP",
  "tags":["mxshop","imooc","bobby","user","web"],
  "user-srv": {
    "host": "本机IP",
    "port": 50051,
    "name": "user-srv"
  },
  "jwt": {
    "key": "jwt认证key"
  },
  "sms": {
    "key": "阿里短信的key",
    "secrect": "阿里短信的secrect",
    "expire": 300
  },
  "redis": {
    "host": "121.40.213.174",
    "port": 6301
  },
  "consul": {
    "host": "服务器（虚拟机）IP",
    "port": 8500
  }
}
```

## userop-web.json
Group：dev
```
{
  "host":"本机IP",
  "tags":["mxshop","imooc","bobby","userop","web"],
  "name": "userop-web",
  "port": 8027,
  "goods-srv": {
    "name": "goods-srv"
  },
  "userop-srv": {
    "name": "userop-srv"
  },
  "jwt": {
    "key": "jwt认证key"
  },
  "consul": {
    "host": "服务器（虚拟机）IP",
    "port": 8500
  }
}
```

## order-web.json
Group：dev
```
{
  "host":"本机IP",
  "tags":["mxshop","imooc","bobby","goods","web"],
  "name": "goods-web",
  "port": 8023,
  "goods-srv": {
    "name": "goods-srv"
  },
  "inventory-srv": {
    "name": "inventory-srv"
  },
  "jwt": {
    "key": "jwt认证key"
  },
  "consul": {
    "host": "服务器（虚拟机）IP",
    "port": 8500
  }
}
```

## order-web.json
Group：dev
```
{
  "host":"本机IP",
  "tags":["mxshop","imooc","bobby","oss","web"],
  "name": "oss-web",
  "port": 8029,
  "oss": {
    "key":"阿里云oss存储的key",
    "secrect":"阿里云oss存储的secrect",
    "host":"阿里云oss存储的Bucket域名",
    "callback_url":"上传后回调",
    "upload_dir":"mxshop-images/"
  },
  "jwt": {
    "key": "jwt认证key"
  },
  "consul": {
    "host": "服务器（虚拟机）IP",
    "port": 8500
  }
}
```

## order-web.json
Group：dev
```
{
  "host":"本机IP",
  "tags":["mxshop","imooc","bobby","order","web"],
  "name": "order-web",
  "port": 8024,
  "goods-srv": {
    "name": "goods-srv"
  },
  "order-srv": {
    "name": "order-srv"
  },
  "inventory-srv": {
    "name": "inventory-srv"
  },
  "jwt": {
    "key": "jwt认证key"
  },
  "consul": {
    "host": "服务器（虚拟机）IP",
    "port": 8500
  },
  "alipay": {
    "app_id":"支付宝支付的id",
    "private_key":"自己的key",
    "ali_public_key":"支付宝的key",
    "notify_url":"支付成功后回调url",
    "return_url":"支付后跳转url"
  }
}
```
