## user-srv-golang.json
Group：dev
```
{
    "name":"user-srv",
    "host":"自己主机的IP",
    "tags":["1","dercy","go","srv"],
    "mysql":{
        "db":"shop_user_srv2",
        "host":"localhost",
        "port":3306,
        "user":"root",
        "password":"123456"
    },
    "consul":{
        "host":"服务器(虚拟机)的IP",
        "port":8500
    }
}
```

## userop-srv-golang.json
Group：dev
```
{
    "name":"userop-srv",
    "host":"自己主机的IP",
    "tags":["1","dercy","golang","userop","srv"],
    "mysql":{
        "db":"shop_userop_srv2",
        "host":"localhost",
        "port":3306,
        "user":"root",
        "password":"123456"
    },
    "consul":{
        "host":"服务器(虚拟机)的IP",
        "port":8500
    }
}
```

## goods-srv-golang.json
Group：dev
```
{
    "name":"goods-srv",
    "host":"自己主机的IP",
    "tags":["1","dercy","goods","srv"],
    "mysql":{
        "db":"shop_goods_srv2",
        "host":"localhost",
        "port":3306,
        "user":"root",
        "password":"123456"
    },
    "consul":{
        "host":"服务器(虚拟机)的IP",
        "port":8500
    },
    "es":{
        "host":"服务器(虚拟机)的IP",
        "port":9200
    }
}
```

## order-srv-golang.json
Group：dev
```
{
    "name":"order-srv",
    "host":"自己主机的IP",
    "tags":["1","dercy","go","order","srv"],
    "mysql":{
        "db":"shop_order_srv2",
        "host":"localhost",
        "port":3306,
        "user":"root",
        "password":"123456"
    },
    "goods_srv":{
        "name":"goods-srv"
    },
    "inventory_srv":{
        "name":"inventory-srv"
    },
    "consul":{
        "host":"服务器(虚拟机)的IP",
        "port":8500
    },
    "jaeger": {
      "host": "192.168.10.130",
      "port": 6831,
      "name": "shop"
    }
}
```