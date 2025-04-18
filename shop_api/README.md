# 电商系统 - web层

API测试：https://nilszhai.apifox.cn

向外提供http接口，服务之间互相解耦

## 项目架构：
- goods-web（商品服务）
- oss-web（oss存储服务）
- order-web（订单和购物车服务）
- user-web（用户服务）
- userop-web（用户操作服务）
- runTest.bat（批量启动）
### 子目录：
- api（核心逻辑）
- config（配置信息-nacos读取）
- forms（表单认证）数据库
- global（公用中间件）
- initialze（初始化中间件）
- models（公用模型[jwt]）
- proto（protobuf混合语言数据标准）
- router（子路由）
- config-*.yaml（nacos配置文件）
- tests（测试用例）
- utils（公用第三方函数[consul]）
- validator（验证器[手机号]）
- main.go（主函数）
- outher.txt（go生成proto）

## nacos详细配置：
https://gitee.com/dercy/shop/blob/master/srv_golang/nacosInfo.md

## 用户服务：
>密码是md5盐值加密，认证方式是JWT认证

>文章：https://blog.csdn.net/the_shy_faker/article/details/127773564

>解决前后端的跨域问题（middlewares/cors.go）

1. 用户列表接口
2. 登录接口
3. 注册接口
5. 阿里云发送短信接口
6. 获取图片验证码

## 商品服务：

>查询商品srv层采用了Elasticsearch全文搜索

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

## 订单和购物车服务
>模仿了jd的逻辑，新建订单是从购物车获取已勾选的商品
1. 购物车列表
2. 添加商品到购物车接口
3. 更新购物车、删除购物车记录接口
4. 订单列表页接口
5. 查询订单详情接口
6. 新建订单接口
### 新建订单后生成的支付功能
    有超时时间 - 过期后归还库存
    1.返回alipay_url支付链接
    2.回调功能

## 用户操作服务
1. 用户添加收藏、取消收藏、查询收藏状态
2. 用户添加地址、更新地址、删除地址、地址列表
3. 用户添加留言、留言列表



