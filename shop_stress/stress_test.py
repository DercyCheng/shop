#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import subprocess
import time
import json
import argparse
from datetime import datetime
import platform

# 导入配置
try:
    from config import SERVICE_PORTS, API_PATHS, DEFAULT_CONFIG
except ImportError:
    print("❌ 错误: 未找到配置文件 (config.py)，请确保该文件存在")
    sys.exit(1)

class StressTest:
    def __init__(self, args):
        self.duration = args.duration
        self.connections = args.connections
        self.threads = args.threads
        self.service = args.service
        self.output_dir = args.output_dir
        self.host = args.host
        
        # 创建输出目录
        os.makedirs(self.output_dir, exist_ok=True)
        
        # 检查wrk是否安装
        self.check_wrk()
    
    def check_wrk(self):
        """检查wrk是否已安装"""
        try:
            subprocess.run(["wrk", "--version"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            print("✅ wrk 已安装")
        except FileNotFoundError:
            print("❌ 错误: wrk 未安装，请先安装wrk")
            print("Windows: 可以使用WSL或下载Windows版本 (https://github.com/wg/wrk/wiki/Installing-wrk-on-Windows)")
            print("Linux: apt-get install wrk 或 yum install wrk")
            print("MacOS: brew install wrk")
            sys.exit(1)
    
    def create_lua_script(self, api_config):
        """创建Lua脚本用于wrk测试，特别是对于POST请求"""
        if api_config.get("method") == "POST" and api_config.get("payload"):
            script_content = f'''
            wrk.method = "{api_config['method']}"
            wrk.body = '{api_config['payload']}'
            wrk.headers["Content-Type"] = "application/json"
            '''
            
            script_file = os.path.join(self.output_dir, f"request_{api_config['path'].replace('/', '_')}.lua")
            with open(script_file, "w") as f:
                f.write(script_content)
            return script_file
        return None
    
    def run_test(self):
        """运行压力测试"""
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        report_file = os.path.join(self.output_dir, f"report_{self.service}_{timestamp}.json")
        results = []
        
        if self.service not in SERVICE_PORTS and self.service != "all":
            print(f"❌ 错误: 未知服务 '{self.service}'. 可用服务: {', '.join(SERVICE_PORTS.keys())} 或 'all'")
            sys.exit(1)
        
        services_to_test = [self.service] if self.service != "all" else SERVICE_PORTS.keys()
        
        for service in services_to_test:
            if service not in API_PATHS:
                print(f"⚠️ 警告: 服务 '{service}' 没有定义API路径配置，跳过测试")
                continue
                
            port = SERVICE_PORTS[service]
            base_url = f"http://{self.host}:{port}"
            
            print(f"\n🚀 开始对 {service} 服务 ({base_url}) 进行压力测试...")
            
            for api_config in API_PATHS[service]:
                url = f"{base_url}{api_config['path']}"
                method = api_config.get("method", "GET")
                description = api_config.get("description", api_config['path'])
                
                print(f"\n📌 测试API: {description} ({method} {url})")
                
                # 为POST请求创建Lua脚本
                lua_script = self.create_lua_script(api_config)
                
                # 构建wrk命令
                cmd = [
                    "wrk",
                    "-t", str(self.threads),
                    "-c", str(self.connections),
                    "-d", f"{self.duration}s",
                    "--latency"
                ]
                
                if lua_script:
                    cmd.extend(["-s", lua_script])
                
                cmd.append(url)
                
                # 运行测试
                try:
                    print(f"⏳ 运行命令: {' '.join(cmd)}")
                    process = subprocess.run(cmd, capture_output=True, text=True)
                    output = process.stdout
                    
                    # 解析结果
                    result = self.parse_wrk_output(output)
                    result["service"] = service
                    result["api"] = api_config['path']
                    result["method"] = method
                    result["description"] = description
                    result["url"] = url
                    result["timestamp"] = datetime.now().isoformat()
                    
                    results.append(result)
                    
                    # 打印结果摘要
                    print(f"✅ 测试完成: {result['requests_per_sec']:.2f} 请求/秒, "
                          f"平均延迟: {result['latency_avg']:.2f}ms, "
                          f"错误率: {result.get('errors_percent', 0):.2f}%")
                    
                except Exception as e:
                    print(f"❌ 测试失败: {str(e)}")
        
        # 保存完整结果到JSON文件
        with open(report_file, 'w', encoding='utf-8') as f:
            json.dump(results, f, ensure_ascii=False, indent=2)
        
        print(f"\n📊 测试报告已保存到: {report_file}")
        self.generate_report(results, report_file.replace('.json', '.html'))
    
    def parse_wrk_output(self, output):
        """解析wrk工具的输出"""
        result = {}
        
        try:
            # 提取请求/秒
            rps_line = [line for line in output.split('\n') if "Requests/sec:" in line]
            if rps_line:
                result["requests_per_sec"] = float(rps_line[0].split(':')[1].strip())
            
            # 提取延迟信息
            latency_lines = output.split("Latency Distribution")[1].split("\n") if "Latency Distribution" in output else []
            for line in latency_lines:
                if "50%" in line:
                    result["latency_50th"] = float(line.strip().split()[1].replace("ms", ""))
                elif "75%" in line:
                    result["latency_75th"] = float(line.strip().split()[1].replace("ms", ""))
                elif "90%" in line:
                    result["latency_90th"] = float(line.strip().split()[1].replace("ms", ""))
                elif "99%" in line:
                    result["latency_99th"] = float(line.strip().split()[1].replace("ms", ""))
            
            # 提取平均延迟和其他统计信息
            threads_line = [line for line in output.split('\n') if "Thread Stats" in line]
            if threads_line:
                stats_line = output.split('\n')[output.split('\n').index(threads_line[0]) + 1]
                parts = stats_line.split()
                result["latency_avg"] = float(parts[1].replace("ms", ""))
                result["latency_stdev"] = float(parts[2].replace("ms", ""))
                result["latency_max"] = float(parts[3].replace("ms", ""))
            
            # 提取传输速率
            transfer_line = [line for line in output.split('\n') if "Transfer/sec:" in line]
            if transfer_line:
                result["transfer_per_sec"] = transfer_line[0].split(':')[1].strip()
            
            # 提取请求总数和错误数
            requests_line = [line for line in output.split('\n') if "requests in" in line][0]
            result["total_requests"] = int(requests_line.split()[0])
            
            # 检查是否有错误
            if "Non-2xx or 3xx responses:" in output:
                error_line = [line for line in output.split('\n') if "Non-2xx or 3xx responses:" in line][0]
                result["error_responses"] = int(error_line.split(':')[1].strip())
                result["errors_percent"] = (result["error_responses"] / result["total_requests"]) * 100
            else:
                result["error_responses"] = 0
                result["errors_percent"] = 0
                
        except Exception as e:
            print(f"解析wrk输出时出错: {str(e)}")
            print(f"原始输出: {output}")
        
        return result
    
    def generate_report(self, results, html_file):
        """生成HTML格式的测试报告"""
        html_content = """
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <title>MXShop 服务压力测试报告</title>
            <style>
                body { font-family: Arial, sans-serif; margin: 20px; }
                h1, h2 { color: #333; }
                table { border-collapse: collapse; width: 100%; margin-bottom: 20px; }
                th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
                th { background-color: #f2f2f2; }
                tr:hover { background-color: #f5f5f5; }
                .container { max-width: 1200px; margin: 0 auto; }
                .summary { padding: 15px; background-color: #f9f9f9; border-radius: 5px; margin-bottom: 20px; }
                .chart-container { height: 400px; margin-bottom: 30px; }
                .error { color: red; }
                .success { color: green; }
            </style>
            <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
        </head>
        <body>
            <div class="container">
                <h1>MXShop 服务压力测试报告</h1>
                <div class="summary">
                    <p><strong>测试时间:</strong> """ + datetime.now().strftime("%Y-%m-%d %H:%M:%S") + """</p>
                    <p><strong>测试配置:</strong> 持续时间 """ + str(self.duration) + """秒, """ + str(self.connections) + """ 连接, """ + str(self.threads) + """ 线程</p>
                </div>
                
                <h2>测试结果摘要</h2>
                <table>
                    <tr>
                        <th>服务</th>
                        <th>API</th>
                        <th>方法</th>
                        <th>请求/秒</th>
                        <th>平均延迟(ms)</th>
                        <th>P99延迟(ms)</th>
                        <th>错误率(%)</th>
                    </tr>
        """
        
        # 添加结果行
        for result in results:
            error_class = "error" if result.get("errors_percent", 0) > 5 else "success"
            html_content += f"""
                <tr>
                    <td>{result['service']}</td>
                    <td>{result['description']}</td>
                    <td>{result['method']}</td>
                    <td>{result.get('requests_per_sec', 0):.2f}</td>
                    <td>{result.get('latency_avg', 0):.2f}</td>
                    <td>{result.get('latency_99th', 0):.2f}</td>
                    <td class="{error_class}">{result.get('errors_percent', 0):.2f}%</td>
                </tr>
            """
        
        # 创建请求/秒图表的数据
        labels = [f"{r['service']} - {r['description']}" for r in results]
        rps_values = [r.get('requests_per_sec', 0) for r in results]
        latency_values = [r.get('latency_avg', 0) for r in results]
        
        html_content += """
                </table>
                
                <h2>性能图表</h2>
                <div class="chart-container">
                    <canvas id="rpsChart"></canvas>
                </div>
                
                <div class="chart-container">
                    <canvas id="latencyChart"></canvas>
                </div>
                
                <script>
                    // 请求/秒图表
                    const rpsCtx = document.getElementById('rpsChart').getContext('2d');
                    new Chart(rpsCtx, {
                        type: 'bar',
                        data: {
                            labels: """ + json.dumps(labels) + """,
                            datasets: [{
                                label: '请求/秒',
                                data: """ + json.dumps(rps_values) + """,
                                backgroundColor: 'rgba(54, 162, 235, 0.6)',
                                borderColor: 'rgba(54, 162, 235, 1)',
                                borderWidth: 1
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                title: {
                                    display: true,
                                    text: '各API请求/秒 (RPS) 对比'
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    title: {
                                        display: true,
                                        text: '请求/秒'
                                    }
                                }
                            }
                        }
                    });
                    
                    // 延迟图表
                    const latencyCtx = document.getElementById('latencyChart').getContext('2d');
                    new Chart(latencyCtx, {
                        type: 'bar',
                        data: {
                            labels: """ + json.dumps(labels) + """,
                            datasets: [{
                                label: '平均延迟 (ms)',
                                data: """ + json.dumps(latency_values) + """,
                                backgroundColor: 'rgba(255, 99, 132, 0.6)',
                                borderColor: 'rgba(255, 99, 132, 1)',
                                borderWidth: 1
                            }]
                        },
                        options: {
                            responsive: true,
                            maintainAspectRatio: false,
                            plugins: {
                                title: {
                                    display: true,
                                    text: '各API平均延迟对比 (ms)'
                                }
                            },
                            scales: {
                                y: {
                                    beginAtZero: true,
                                    title: {
                                        display: true,
                                        text: '延迟 (ms)'
                                    }
                                }
                            }
                        }
                    });
                </script>
            </div>
        </body>
        </html>
        """
        
        with open(html_file, 'w', encoding='utf-8') as f:
            f.write(html_content)
        
        print(f"📈 HTML报告已生成: {html_file}")


def main():
    parser = argparse.ArgumentParser(description='MXShop服务压力测试工具')
    parser.add_argument('-s', '--service', default='all', 
                        help='要测试的服务 (user, goods, order, userop, oss 或 all)')
    parser.add_argument('-d', '--duration', type=int, default=DEFAULT_CONFIG.get('duration', 10), 
                        help='测试持续时间(秒)')
    parser.add_argument('-c', '--connections', type=int, default=DEFAULT_CONFIG.get('connections', 100), 
                        help='并发连接数')
    parser.add_argument('-t', '--threads', type=int, default=DEFAULT_CONFIG.get('threads', 4), 
                        help='使用的线程数')
    parser.add_argument('-o', '--output-dir', default='./results', 
                        help='测试结果输出目录')
    parser.add_argument('--host', default=DEFAULT_CONFIG.get('host', 'localhost'), 
                        help='服务主机地址')
    
    args = parser.parse_args()
    
    # 创建输出目录
    if not os.path.exists(args.output_dir):
        os.makedirs(args.output_dir)

    stress_test = StressTest(args)
    stress_test.run_test()


if __name__ == '__main__':
    main()