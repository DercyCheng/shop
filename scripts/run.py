#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Shop 电商微服务系统启动脚本
用于管理和启动各个微服务组件
支持Docker、Kubernetes和本地部署
"""

import os
import sys
import subprocess
import argparse
import time
import signal
import platform
import shutil

# 项目根目录
ROOT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# 颜色输出
class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

# 服务配置
SERVICES = {
    # 基础设施服务
    "infrastructure": [
        {"name": "MySQL", "command": "docker-compose up -d mysql", "wait": 10},
        {"name": "Redis", "command": "docker-compose up -d redis", "wait": 5},
        {"name": "Elasticsearch", "command": "docker-compose up -d elasticsearch", "wait": 20},
        {"name": "RocketMQ", "command": "docker-compose up -d rocketmq-namesrv rocketmq-broker", "wait": 10},
        {"name": "Consul", "command": "docker-compose up -d consul", "wait": 5},
        {"name": "Nacos", "command": "docker-compose up -d nacos", "wait": 15},
        {"name": "Jaeger", "command": "docker-compose up -d jaeger", "wait": 5},
        {"name": "Nginx", "command": "docker-compose up -d nginx", "wait": 3},
    ],
    
    # 服务层 (SRV)
    "srv": [
        {"name": "用户服务", "dir": "shop_srv/user_srv", "command": "go run cmd/server/main.go -p 50051"},
        {"name": "商品服务", "dir": "shop_srv/goods_srv", "command": "go run cmd/server/main.go -p 50052"},
        {"name": "库存服务", "dir": "shop_srv/inventory_srv", "command": "go run cmd/server/main.go -p 50053"},
        {"name": "订单服务", "dir": "shop_srv/order_srv", "command": "go run cmd/server/main.go -p 50054"},
        {"name": "用户操作服务", "dir": "shop_srv/userop_srv", "command": "go run cmd/server/main.go -p 50055"}
    ],
    
    # API层 (Web)
    "api": [
        {"name": "用户API", "dir": "shop_web/user-web", "command": "go run main.go"},
        {"name": "商品API", "dir": "shop_web/goods-web", "command": "go run main.go"},
        {"name": "订单API", "dir": "shop_web/order-web", "command": "go run main.go"},
        {"name": "用户操作API", "dir": "shop_web/userop-web", "command": "go run main.go"},
        {"name": "OSS服务API", "dir": "shop_web/oss-web", "command": "go run main.go"}
    ],

    # K8s配置
    "k8s": [
        {"name": "基础设施", "file": "deploy/k8s/infrastructure/infrastructure.yaml"},
        {"name": "服务层", "file": "deploy/k8s/services/services.yaml"},
        {"name": "API层", "file": "deploy/k8s/web/web.yaml"}
    ]
}

# 运行中的进程
processes = {}

def print_banner():
    """打印启动横幅"""
    banner = f"""
{Colors.BLUE}{Colors.BOLD}
  ____  _                   ____            _                 
 / ___|| |__   ___  _ __   / ___|  ___ _ __(__) ___  ___ ___ 
 \\___ \\| '_ \\ / _ \\| '_ \\  \\___ \\ / _ \\ '__| |/ __|/ _ / __|
  ___) | | | | (_) | |_) |  ___) |  __/ |  | | (__|  __\\__ \\
 |____/|_| |_|\\___/| .__/  |____/ \\___|_|  |_|\\___|\\___|___/
                   |_|                                       
{Colors.ENDC}
Shop 电商微服务系统启动工具
"""
    print(banner)

def run_command(command, cwd=None, background=True):
    """运行命令"""
    shell = platform.system() == "Windows"
    if background:
        if platform.system() == "Windows":
            # Windows下启动单独的命令行窗口
            return subprocess.Popen(
                command, 
                cwd=cwd, 
                shell=shell, 
                creationflags=subprocess.CREATE_NEW_CONSOLE
            )
        else:
            # Unix/Linux/Mac下使用nohup实现后台运行
            # 重定向输出到临时文件
            log_file = os.path.join(cwd if cwd else os.getcwd(), "service.log")
            command = f"nohup {command} > {log_file} 2>&1 &"
            process = subprocess.Popen(
                command,
                cwd=cwd,
                shell=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            # 等待nohup启动
            time.sleep(1)
            # 返回实际运行的子进程ID（不是当前Python进程的子进程）
            # 注意这是近似值，对于进程管理不够完善，但对于演示目的足够
            try:
                if cwd is None:
                    cwd = os.getcwd()
                # 找到刚才通过nohup启动的进程
                find_cmd = f"ps -ef | grep '{command.replace('&', '')}' | grep -v grep | awk '{{print $2}}'"
                pid = subprocess.check_output(find_cmd, shell=True).decode().strip()
                if pid:
                    return pid  # 返回进程ID而不是Popen对象
            except Exception as e:
                print(f"{Colors.YELLOW}警告: 无法获取后台进程ID: {e}{Colors.ENDC}")
            return process
    else:
        # 前台运行
        return subprocess.run(command, cwd=cwd, shell=shell)

def start_infrastructure():
    """启动基础设施服务"""
    print(f"{Colors.HEADER}正在启动基础设施服务...{Colors.ENDC}")
    
    for service in SERVICES["infrastructure"]:
        print(f"{Colors.BLUE}启动 {service['name']}...{Colors.ENDC}")
        try:
            run_command(service["command"], cwd=ROOT_DIR, background=False)
            print(f"{Colors.GREEN}{service['name']} 启动命令已执行{Colors.ENDC}")
            print(f"{Colors.YELLOW}等待 {service['wait']} 秒钟让服务完全启动...{Colors.ENDC}")
            time.sleep(service["wait"])
        except Exception as e:
            print(f"{Colors.RED}启动 {service['name']} 失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}基础设施服务已启动{Colors.ENDC}")

def start_srv_services():
    """启动服务层"""
    print(f"{Colors.HEADER}正在启动服务层 (SRV)...{Colors.ENDC}")
    
    for service in SERVICES["srv"]:
        print(f"{Colors.BLUE}启动 {service['name']}...{Colors.ENDC}")
        try:
            service_dir = os.path.join(ROOT_DIR, service["dir"])
            process = run_command(service["command"], cwd=service_dir)
            processes[service["name"]] = process
            print(f"{Colors.GREEN}{service['name']} 已启动{Colors.ENDC}")
            # 给服务一些启动时间
            time.sleep(3)
        except Exception as e:
            print(f"{Colors.RED}启动 {service['name']} 失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}服务层 (SRV) 已启动{Colors.ENDC}")

def start_api_services():
    """启动API层"""
    print(f"{Colors.HEADER}正在启动API层 (Web)...{Colors.ENDC}")
    
    for service in SERVICES["api"]:
        print(f"{Colors.BLUE}启动 {service['name']}...{Colors.ENDC}")
        try:
            service_dir = os.path.join(ROOT_DIR, service["dir"])
            process = run_command(service["command"], cwd=service_dir)
            processes[service["name"]] = process
            print(f"{Colors.GREEN}{service['name']} 已启动{Colors.ENDC}")
            # 给服务一些启动时间
            time.sleep(2)
        except Exception as e:
            print(f"{Colors.RED}启动 {service['name']} 失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}API层 (Web) 已启动{Colors.ENDC}")

def init_database():
    """初始化数据库"""
    print(f"{Colors.HEADER}正在初始化数据库...{Colors.ENDC}")
    
    try:
        sql_file = os.path.join(ROOT_DIR, "scripts", "init.sql")
        # 使用docker中的mysql客户端执行SQL文件
        command = f"docker exec -i $(docker ps -qf 'name=mysql') mysql -uroot -proot < {sql_file}"
        run_command(command, background=False)
        print(f"{Colors.GREEN}数据库初始化完成{Colors.ENDC}")
    except Exception as e:
        print(f"{Colors.RED}数据库初始化失败: {e}{Colors.ENDC}")

def stop_services():
    """停止所有服务"""
    print(f"{Colors.HEADER}正在停止所有服务...{Colors.ENDC}")
    
    # 停止API和SRV层服务
    for name, process in processes.items():
        print(f"{Colors.BLUE}停止 {name}...{Colors.ENDC}")
        try:
            if isinstance(process, str):  # 如果是进程ID
                os.kill(int(process), signal.SIGTERM)
            else:
                process.terminate()
                process.wait(timeout=5)
            print(f"{Colors.GREEN}{name} 已停止{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}停止 {name} 失败: {e}{Colors.ENDC}")
            # 强制终止
            try:
                if isinstance(process, str):
                    os.kill(int(process), signal.SIGKILL)
                else:
                    process.kill()
            except:
                pass
    
    # 停止基础设施服务
    print(f"{Colors.BLUE}停止基础设施服务...{Colors.ENDC}")
    try:
        run_command("docker-compose down", cwd=ROOT_DIR, background=False)
        print(f"{Colors.GREEN}基础设施服务已停止{Colors.ENDC}")
    except Exception as e:
        print(f"{Colors.RED}停止基础设施服务失败: {e}{Colors.ENDC}")
    
    processes.clear()
    print(f"{Colors.GREEN}所有服务已停止{Colors.ENDC}")

def show_status():
    """显示服务状态"""
    print(f"{Colors.HEADER}服务状态:{Colors.ENDC}")
    
    # 显示Docker服务状态
    print(f"{Colors.BLUE}基础设施服务:{Colors.ENDC}")
    try:
        run_command("docker-compose ps", cwd=ROOT_DIR, background=False)
    except Exception as e:
        print(f"{Colors.RED}获取基础设施服务状态失败: {e}{Colors.ENDC}")
    
    # 显示微服务状态
    print(f"\n{Colors.BLUE}微服务状态:{Colors.ENDC}")
    if not processes:
        print(f"{Colors.YELLOW}没有运行中的微服务{Colors.ENDC}")
    else:
        for name in processes:
            print(f"- {name}: {Colors.GREEN}运行中{Colors.ENDC}")

def build_docker_images():
    """构建Docker镜像"""
    print(f"{Colors.HEADER}正在构建Docker镜像...{Colors.ENDC}")
    
    # 构建服务层镜像
    print(f"{Colors.BLUE}构建服务层 (SRV) 镜像...{Colors.ENDC}")
    srv_services = [
        {"path": "shop_srv/user_srv", "name": "user-srv", "main": "cmd/server/main.go"},
        {"path": "shop_srv/goods_srv", "name": "goods-srv", "main": "cmd/server/main.go"},
        {"path": "shop_srv/inventory_srv", "name": "inventory-srv", "main": "cmd/server/main.go"},
        {"path": "shop_srv/order_srv", "name": "order-srv", "main": "cmd/server/main.go"},
        {"path": "shop_srv/userop_srv", "name": "userop-srv", "main": "cmd/server/main.go"}
    ]
    
    for service in srv_services:
        print(f"{Colors.BLUE}构建 {service['name']} 镜像...{Colors.ENDC}")
        try:
            service_path = os.path.join(ROOT_DIR, service["path"])
            command = f"docker build -t shop/{service['name']}:latest -f {ROOT_DIR}/Dockerfile --build-arg SERVICE_PATH={service['main']} --build-arg SERVICE_NAME={service['name']} {service_path}"
            run_command(command, background=False)
            print(f"{Colors.GREEN}{service['name']} 镜像构建完成{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}构建 {service['name']} 镜像失败: {e}{Colors.ENDC}")
    
    # 构建API层镜像
    print(f"{Colors.BLUE}构建API层 (Web) 镜像...{Colors.ENDC}")
    web_services = [
        {"path": "shop_web/user-web", "name": "user-web", "main": "main.go"},
        {"path": "shop_web/goods-web", "name": "goods-web", "main": "main.go"},
        {"path": "shop_web/order-web", "name": "order-web", "main": "main.go"},
        {"path": "shop_web/userop-web", "name": "userop-web", "main": "main.go"},
        {"path": "shop_web/oss-web", "name": "oss-web", "main": "main.go"}
    ]
    
    for service in web_services:
        print(f"{Colors.BLUE}构建 {service['name']} 镜像...{Colors.ENDC}")
        try:
            service_path = os.path.join(ROOT_DIR, service["path"])
            command = f"docker build -t shop/{service['name']}:latest -f {ROOT_DIR}/Dockerfile --build-arg SERVICE_PATH={service['main']} --build-arg SERVICE_NAME={service['name']} {service_path}"
            run_command(command, background=False)
            print(f"{Colors.GREEN}{service['name']} 镜像构建完成{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}构建 {service['name']} 镜像失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}所有Docker镜像构建完成{Colors.ENDC}")

def deploy_to_k8s():
    """部署到Kubernetes"""
    print(f"{Colors.HEADER}正在部署到Kubernetes...{Colors.ENDC}")
    
    # 检查是否有kubectl
    try:
        run_command("kubectl version --client", background=False)
    except Exception as e:
        print(f"{Colors.RED}未找到kubectl，请确保安装了Kubernetes客户端工具: {e}{Colors.ENDC}")
        return
    
    # 创建命名空间
    print(f"{Colors.BLUE}创建Kubernetes命名空间...{Colors.ENDC}")
    try:
        run_command("kubectl create namespace shop-system --dry-run=client -o yaml | kubectl apply -f -", background=False)
        print(f"{Colors.GREEN}命名空间已创建{Colors.ENDC}")
    except Exception as e:
        print(f"{Colors.YELLOW}创建命名空间可能已存在: {e}{Colors.ENDC}")
    
    # 准备ConfigMap数据
    print(f"{Colors.BLUE}准备配置文件...{Colors.ENDC}")
    try:
        # 处理MySQL初始化脚本
        sql_file = os.path.join(ROOT_DIR, "scripts", "init.sql")
        if os.path.exists(sql_file):
            with open(sql_file, 'r') as f:
                sql_content = f.read()
            # 更新MySQL init ConfigMap
            command = f"kubectl create configmap mysql-init-scripts --from-literal=init.sql='{sql_content}' -n shop-system --dry-run=client -o yaml | kubectl apply -f -"
            run_command(command, background=False)
            print(f"{Colors.GREEN}MySQL初始化脚本已准备{Colors.ENDC}")
    except Exception as e:
        print(f"{Colors.YELLOW}准备MySQL初始化脚本失败，将使用默认值: {e}{Colors.ENDC}")
    
    # 部署各个组件
    for component in SERVICES["k8s"]:
        print(f"{Colors.BLUE}部署 {component['name']}...{Colors.ENDC}")
        try:
            component_file = os.path.join(ROOT_DIR, component["file"])
            command = f"kubectl apply -f {component_file}"
            run_command(command, background=False)
            print(f"{Colors.GREEN}{component['name']} 已部署{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}部署 {component['name']} 失败: {e}{Colors.ENDC}")
    
    # 等待服务启动
    print(f"{Colors.YELLOW}等待服务启动...{Colors.ENDC}")
    time.sleep(10)
    
    # 检查部署状态
    print(f"{Colors.BLUE}检查部署状态...{Colors.ENDC}")
    try:
        run_command("kubectl get pods -n shop-system", background=False)
        run_command("kubectl get svc -n shop-system", background=False)
    except Exception as e:
        print(f"{Colors.RED}检查部署状态失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}Kubernetes部署完成{Colors.ENDC}")
    
    # 获取API网关地址
    print(f"{Colors.BLUE}获取API网关访问地址...{Colors.ENDC}")
    try:
        run_command("kubectl get svc shop-gateway -n shop-system", background=False)
    except Exception as e:
        print(f"{Colors.RED}获取API网关地址失败: {e}{Colors.ENDC}")

def undeploy_from_k8s():
    """从Kubernetes中卸载"""
    print(f"{Colors.HEADER}正在从Kubernetes卸载...{Colors.ENDC}")
    
    # 检查是否有kubectl
    try:
        run_command("kubectl version --client", background=False)
    except Exception as e:
        print(f"{Colors.RED}未找到kubectl，请确保安装了Kubernetes客户端工具: {e}{Colors.ENDC}")
        return
    
    # 删除各个组件
    for component in reversed(SERVICES["k8s"]):
        print(f"{Colors.BLUE}删除 {component['name']}...{Colors.ENDC}")
        try:
            component_file = os.path.join(ROOT_DIR, component["file"])
            command = f"kubectl delete -f {component_file}"
            run_command(command, background=False)
            print(f"{Colors.GREEN}{component['name']} 已删除{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}删除 {component['name']} 失败: {e}{Colors.ENDC}")
    
    # 询问是否删除命名空间和持久数据
    delete_namespace = input(f"{Colors.YELLOW}是否删除整个命名空间和所有数据? (y/N): {Colors.ENDC}")
    if delete_namespace.lower() == 'y':
        try:
            run_command("kubectl delete namespace shop-system", background=False)
            print(f"{Colors.GREEN}命名空间和所有数据已删除{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}删除命名空间失败: {e}{Colors.ENDC}")
    
    print(f"{Colors.GREEN}Kubernetes卸载完成{Colors.ENDC}")

def show_help():
    """显示帮助信息"""
    print(f"{Colors.HEADER}Shop 电商微服务系统启动工具{Colors.ENDC}")
    print(f"\n{Colors.BOLD}可用命令:{Colors.ENDC}")
    
    # 本地开发命令
    print(f"\n{Colors.BOLD}本地开发命令:{Colors.ENDC}")
    print(f"  {Colors.GREEN}all{Colors.ENDC} - 启动所有服务")
    print(f"  {Colors.GREEN}infra{Colors.ENDC} - 只启动基础设施服务")
    print(f"  {Colors.GREEN}srv{Colors.ENDC} - 只启动服务层")
    print(f"  {Colors.GREEN}api{Colors.ENDC} - 只启动API层")
    print(f"  {Colors.GREEN}stop{Colors.ENDC} - 停止所有服务")
    print(f"  {Colors.GREEN}status{Colors.ENDC} - 显示服务状态")
    print(f"  {Colors.GREEN}init-db{Colors.ENDC} - 初始化数据库")
    
    # Docker命令
    print(f"\n{Colors.BOLD}Docker命令:{Colors.ENDC}")
    print(f"  {Colors.GREEN}docker-build{Colors.ENDC} - 构建Docker镜像")
    print(f"  {Colors.GREEN}docker-up{Colors.ENDC} - 使用Docker Compose启动所有服务")
    print(f"  {Colors.GREEN}docker-down{Colors.ENDC} - 使用Docker Compose停止所有服务")
    
    # Kubernetes命令
    print(f"\n{Colors.BOLD}Kubernetes命令:{Colors.ENDC}")
    print(f"  {Colors.GREEN}k8s-deploy{Colors.ENDC} - 部署到Kubernetes")
    print(f"  {Colors.GREEN}k8s-undeploy{Colors.ENDC} - 从Kubernetes卸载")
    
    # 帮助信息
    print(f"\n{Colors.BOLD}其他命令:{Colors.ENDC}")
    print(f"  {Colors.GREEN}help{Colors.ENDC} - 显示此帮助信息")
    
    print(f"\n{Colors.BOLD}示例:{Colors.ENDC}")
    print(f"  python run.py all    # 启动所有服务")
    print(f"  python run.py docker-build   # 构建Docker镜像")
    print(f"  python run.py k8s-deploy   # 部署到Kubernetes")

def main():
    """主函数"""
    print_banner()
    
    parser = argparse.ArgumentParser(description="Shop 电商微服务系统启动工具")
    parser.add_argument("command", nargs="?", default="help", 
                        choices=["all", "infra", "srv", "api", "stop", "status", "init-db", 
                                "docker-build", "docker-up", "docker-down", 
                                "k8s-deploy", "k8s-undeploy", "help"],
                        help="要执行的命令")
    
    args = parser.parse_args()
    
    if args.command == "all":
        start_infrastructure()
        init_database()
        start_srv_services()
        start_api_services()
        print(f"\n{Colors.GREEN}{Colors.BOLD}所有服务已启动!{Colors.ENDC}")
        print(f"\n{Colors.BLUE}访问服务:{Colors.ENDC}")
        print(f"  - Swagger API文档: http://localhost:8021/swagger/index.html")
        print(f"  - Nacos控制台: http://localhost:8848/nacos (user/pass: nacos/nacos)")
        print(f"  - Consul UI: http://localhost:8500")
        print(f"  - Jaeger UI: http://localhost:16686")
        print(f"  - API网关: http://localhost")
    
    elif args.command == "infra":
        start_infrastructure()
    
    elif args.command == "srv":
        start_srv_services()
    
    elif args.command == "api":
        start_api_services()
    
    elif args.command == "stop":
        stop_services()
    
    elif args.command == "status":
        show_status()
    
    elif args.command == "init-db":
        init_database()
    
    elif args.command == "docker-build":
        build_docker_images()
    
    elif args.command == "docker-up":
        print(f"{Colors.BLUE}使用Docker Compose启动所有服务...{Colors.ENDC}")
        try:
            run_command("docker-compose up -d", cwd=ROOT_DIR, background=False)
            print(f"{Colors.GREEN}所有服务已启动!{Colors.ENDC}")
            print(f"\n{Colors.BLUE}访问服务:{Colors.ENDC}")
            print(f"  - API网关: http://localhost")
            print(f"  - Nacos控制台: http://localhost:8848/nacos (user/pass: nacos/nacos)")
            print(f"  - Consul UI: http://localhost:8500")
            print(f"  - Jaeger UI: http://localhost:16686")
        except Exception as e:
            print(f"{Colors.RED}启动服务失败: {e}{Colors.ENDC}")
    
    elif args.command == "docker-down":
        print(f"{Colors.BLUE}使用Docker Compose停止所有服务...{Colors.ENDC}")
        try:
            run_command("docker-compose down", cwd=ROOT_DIR, background=False)
            print(f"{Colors.GREEN}所有服务已停止{Colors.ENDC}")
        except Exception as e:
            print(f"{Colors.RED}停止服务失败: {e}{Colors.ENDC}")
    
    elif args.command == "k8s-deploy":
        # 检查是否已经构建了镜像
        print(f"{Colors.YELLOW}确保已经构建了Docker镜像并推送到可访问的注册表{Colors.ENDC}")
        confirm = input(f"{Colors.YELLOW}是否继续部署? (y/N): {Colors.ENDC}")
        if confirm.lower() == 'y':
            deploy_to_k8s()
    
    elif args.command == "k8s-undeploy":
        undeploy_from_k8s()
    
    elif args.command == "help":
        show_help()

if __name__ == "__main__":
    main()