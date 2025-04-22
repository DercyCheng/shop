# 电商系统 - server层

API测试：https://nilszhai.apifox.cn

服务之间GRPC调用，提供GRPC进行连接

## 项目架构：
- goods_srv（商品服务）
- inventory_srv（库存服务）
- order_srv（订单和购物车服务）
- user_srv（用户服务）
- userop_srv（用户操作服务）
- tools（生成grpc文件脚本）
- runTest.bat（批量启动）
### 子目录：
- config（配置信息-nacos读取）
- global（公用中间件）
- handler（核心逻辑）
- initialze（初始化中间件）
- model（gorm模型）
- proto（protobuf混合语言数据标准）
- config-debug.yaml（nacos配置文件）
- tests（测试用例）
- main.go（主函数）

## nacos详细配置：
https://gitee.com/jzins/mxshop/blob/master/srv_golang/nacosInfo.md

## 用户服务：
>密码是md5盐值加密，认证方式是JWT认证

>详细参考我的文章：https://blog.csdn.net/the_shy_faker/article/details/127773564

1. 用户列表接口
2. 通过id和mobile查询用户
3. 新建用户
4. 修改用户和校验密码接口

## 商品服务：

>查询商品采用了Elasticsearch全文搜索

1. 品牌列表
2. 品牌新建，删除、更新
3. 轮播图的查询、新增、删除和修改
4. 商品分类的列表接口
5. 获取商品分类的子分类
6. 商品分类的新建，删除和更新接口
7. 品牌分类相关接口
8. 商品列表页接口
9. 批量获取商品信息、商品详情接口
10. 新增、修改和删除商品接口

## 库存服务：
> 使用redsync实现redis分布式锁，解决了库存的互斥性、原子性、安全性、宕机等问题
这里有一个`续租`的坑，世界上没有最好的锁！

>详细参考我的文章：https://blog.csdn.net/the_shy_faker/article/details/127981144
1. 获取库存详情接口
2. 归还库存接口
3. 修改库存接口

## 订单和购物车服务
>模仿了jd的逻辑，新建订单是从购物车获取已勾选的商品
1. 购物车列表
2. 添加商品到购物车接口
3. 更新购物车、删除购物车记录接口
4. 订单列表页接口
5. 查询订单详情接口
6. 新建订单接口

## 用户操作服务
1. 用户添加收藏、取消收藏、查询收藏状态
2. 用户添加地址、更新地址、删除地址、地址列表
3. 用户添加留言、留言列表

