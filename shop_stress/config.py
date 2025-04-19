#!/usr/bin/env python
# -*- coding: utf-8 -*-

# 这是一个配置文件，用于自定义压力测试参数和API路径
# 您可以根据实际API路径和需求修改此文件

# 服务端口映射
SERVICE_PORTS = {
    "user": 8021,
    "goods": 8022,
    "order": 8023,
    "userop": 8024,
    "oss": 8025
}

# 测试路径配置
# 每个服务包含多个API端点的配置
# 每个API配置包含：
# - path: API路径
# - method: HTTP方法 (GET, POST, PUT, DELETE等)
# - description: API描述 
# - payload: POST请求的JSON负载 (仅POST请求需要)
API_PATHS = {
    "user": [
        {"path": "/v1/user/list", "method": "GET", "description": "用户列表"},
        {"path": "/v1/user/login", "method": "POST", "description": "用户登录", 
         "payload": '{"mobile": "18888888888", "password": "admin123"}'}
    ],
    "goods": [
        {"path": "/v1/goods", "method": "GET", "description": "商品列表"},
        {"path": "/v1/goods/1", "method": "GET", "description": "获取商品详情"},
        {"path": "/v1/categories", "method": "GET", "description": "商品类别列表"},
        {"path": "/v1/brands", "method": "GET", "description": "品牌列表"}
    ],
    "order": [
        {"path": "/v1/order", "method": "GET", "description": "订单列表"},
        {"path": "/v1/order/create", "method": "POST", "description": "创建订单", 
         "payload": '{"goods_id": 1, "goods_num": 1, "address": "测试地址"}'}
    ],
    "userop": [
        {"path": "/v1/message", "method": "GET", "description": "用户消息列表"},
        {"path": "/v1/favorite", "method": "GET", "description": "用户收藏列表"},
        {"path": "/v1/address", "method": "GET", "description": "用户地址列表"}
    ],
    "oss": [
        {"path": "/oss/upload", "method": "GET", "description": "上传文件页面"},
        {"path": "/health", "method": "GET", "description": "健康检查"}
    ]
}

# 默认测试配置
DEFAULT_CONFIG = {
    "duration": 10,       # 测试持续时间(秒)
    "connections": 100,   # 并发连接数
    "threads": 4,         # 线程数
    "host": "localhost"   # 默认主机
}