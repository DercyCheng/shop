# MXShop 服务压力测试工具

这是一个针对MXShop微服务架构的压力测试工具，使用wrk作为底层压测引擎，通过Python脚本进行管理和报告生成。

## 功能特点

- 支持对MXShop各微服务的压力测试
- 可自定义测试参数：并发连接数、线程数和测试持续时间
- 支持GET和POST请求测试
- 自动生成JSON和HTML格式的详细测试报告
- 测试报告包含请求/秒、响应时间、错误率等关键指标
- 可视化图表展示测试结果

## 前置要求

- Python 3.6+
- wrk 压力测试工具
  - Linux: `apt-get install wrk` 或 `yum install wrk`
  - MacOS: `brew install wrk`
  - Windows: 可以使用WSL或下载Windows版本 (参考 [安装指南](https://github.com/wg/wrk/wiki/Installing-wrk-on-Windows))

## 安装

1. 确保您已安装Python 3.6或更高版本
2. 安装wrk工具（如上所述）
3. 将`shop_stress`目录克隆或复制到您的MXShop项目中

## 使用方法

### 基本用法

```bash
# 测试所有服务
python shop_stress/stress_test.py

# 测试特定服务
python shop_stress/stress_test.py -s user

# 指定测试参数
python shop_stress/stress_test.py -s goods -d 30 -c 200 -t 8
```

### 命令行参数

- `-s`, `--service`: 要测试的服务名称（user、goods、order、userop、oss或all）
- `-d`, `--duration`: 测试持续时间（秒）
- `-c`, `--connections`: 并发连接数
- `-t`, `--threads`: 测试线程数
- `-o`, `--output-dir`: 测试结果输出目录
- `--host`: 服务主机地址（默认为localhost）

### 示例

```bash
# 测试用户服务 30秒，200并发，8线程
python shop_stress/stress_test.py -s user -d 30 -c 200 -t 8

# 测试商品服务 60秒，500并发，16线程
python shop_stress/stress_test.py -s goods -d 60 -c 500 -t 16 --host 192.168.1.100

# 测试所有服务，将结果输出到指定目录
python shop_stress/stress_test.py -o ./stress_reports
```

## 自定义测试

您可以通过编辑`config.py`文件来自定义测试配置：

1. 修改`SERVICE_PORTS`来更新服务端口
2. 修改`API_PATHS`来添加或更新测试的API端点
3. 修改`DEFAULT_CONFIG`来更改默认测试参数

### 添加新的测试API

编辑`config.py`文件中的`API_PATHS`字典：

```python
API_PATHS = {
    "user": [
        # 现有API...
        {"path": "/v1/user/register", "method": "POST", "description": "用户注册", 
         "payload": '{"username": "testuser", "password": "password123", "mobile": "13800138000"}'}
    ],
    # 其他服务...
}
```

## 测试报告

测试完成后，将在输出目录（默认为`./results`）生成两种格式的报告：

1. JSON文件：包含所有原始测试数据
2. HTML文件：包含格式化的测试结果和可视化图表

HTML报告包括：
- 测试概要信息
- 各API的请求/秒、响应时间和错误率
- 请求/秒对比图表
- 平均响应时间对比图表

## 故障排除

1. 如果运行时出现"wrk未安装"错误，请确保已正确安装wrk工具并添加到PATH环境变量中
2. 如果无法连接到服务，请确保服务正在运行，并检查主机地址和端口配置
3. 对于Windows用户，建议使用WSL环境运行测试

## 许可证

MIT