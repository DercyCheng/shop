import subprocess
import os
import sys
import time

def check_docker():
    """Check if Docker is installed and running"""
    try:
        subprocess.run(["docker", "info"], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("Docker is not running or not installed. Please install and start Docker first.")
        return False

def setup_network():
    """Create Docker network for the application if it doesn't exist"""
    try:
        # Check if network exists
        result = subprocess.run(
            ["docker", "network", "ls", "--filter", "name=mxshop-net", "--format", "{{.Name}}"],
            check=True, capture_output=True, text=True
        )
        
        if "mxshop-net" not in result.stdout:
            # Create network
            subprocess.run(["docker", "network", "create", "mxshop-net"], check=True)
            print("Created Docker network: mxshop-net")
        
        return True
    except subprocess.CalledProcessError as e:
        print(f"Failed to setup Docker network: {e}")
        return False

def start_infrastructure():
    """Start necessary infrastructure containers (MySQL, Nacos, etc.)"""
    try:
        # Start MySQL
        mysql_running = subprocess.run(
            ["docker", "ps", "--filter", "name=mxshop-mysql", "--format", "{{.Names}}"],
            check=True, capture_output=True, text=True
        )
        
        if "mxshop-mysql" not in mysql_running.stdout:
            print("Starting MySQL container...")
            subprocess.run([
                "docker", "run", "-d",
                "--name", "mxshop-mysql",
                "--network", "mxshop-net",
                "-p", "3306:3306",
                "-e", "MYSQL_ROOT_PASSWORD=root",
                "-e", "MYSQL_DATABASE=mxshop",
                "-v", "mxshop-mysql-data:/var/lib/mysql",
                "mysql:8.0"
            ], check=True)
            
            # Wait for MySQL to be ready
            print("Waiting for MySQL to be ready...")
            time.sleep(20)
            
            # Initialize database with SQL script
            print("Initializing MySQL database...")
            try:
                with open("./script/init.sql", "rb") as sql_file:
                    subprocess.run([
                        "docker", "exec", "-i", "mxshop-mysql", 
                        "mysql", "-uroot", "-proot"
                    ], input=sql_file.read(), check=True)
            except FileNotFoundError:
                print("Error: init.sql file not found.")
                return False
            except subprocess.CalledProcessError as e:
                print(f"Error executing SQL script: {e}")
                return False
        
        # Start Nacos
        nacos_running = subprocess.run(
            ["docker", "ps", "--filter", "name=mxshop-nacos", "--format", "{{.Names}}"],
            check=True, capture_output=True, text=True
        )
        
        if "mxshop-nacos" not in nacos_running.stdout:
            print("Starting Nacos container...")
            subprocess.run([
                "docker", "run", "-d",
                "--name", "mxshop-nacos",
                "--network", "mxshop-net",
                "-p", "8848:8848",
                "-e", "MODE=standalone",
                "-e", "SPRING_DATASOURCE_PLATFORM=mysql",
                "-e", "MYSQL_SERVICE_HOST=mxshop-mysql",
                "-e", "MYSQL_SERVICE_PORT=3306",
                "-e", "MYSQL_SERVICE_USER=root",
                "-e", "MYSQL_SERVICE_PASSWORD=root",
                "-e", "MYSQL_SERVICE_DB_NAME=nacos",
                "nacos/nacos-server:latest"
            ], check=True)
            
            # Wait for Nacos to be ready
            print("Waiting for Nacos to be ready...")
            time.sleep(20)
        
        return True
    except subprocess.CalledProcessError as e:
        print(f"Failed to start infrastructure: {e}")
        return False
    except Exception as e:
        print(f"Unexpected error starting infrastructure: {e}")
        return False

def build_and_run():
    """Build and run the MXShop application"""
    try:
        # Build Docker image
        print("Building MXShop Docker image...")
        subprocess.run(["docker", "build", "-t", "mxshop-app", "."], check=True)
        
        # Check if container already exists
        container_exists = subprocess.run(
            ["docker", "ps", "-a", "--filter", "name=mxshop-app", "--format", "{{.Names}}"],
            check=True, capture_output=True, text=True
        )
        
        if "mxshop-app" in container_exists.stdout:
            print("Stopping and removing existing container...")
            subprocess.run(["docker", "rm", "-f", "mxshop-app"], check=True)
        
        # Run Docker container
        print("Starting MXShop application container...")
        subprocess.run([
            "docker", "run", "-d",
            "--name", "mxshop-app",
            "--network", "mxshop-net",
            "-p", "8080:8080",
            "-p", "8021:8021",
            "-p", "8022:8022",
            "-p", "8023:8023",
            "-p", "8024:8024",
            "-p", "8025:8025",
            "-e", "DEV_CONFIG=1",
            "mxshop-app"
        ], check=True)
        
        return True
    except subprocess.CalledProcessError as e:
        print(f"Failed to build and run application: {e}")
        return False
    except Exception as e:
        print(f"Unexpected error building and running application: {e}")
        return False

def main():
    """Main function to orchestrate the setup and startup of the application"""
    try:
        print("Starting MXShop deployment...")
        
        # Check if Docker is available
        if not check_docker():
            sys.exit(1)
        
        # Setup Docker network
        if not setup_network():
            sys.exit(1)
        
        # Start infrastructure (MySQL, Nacos)
        if not start_infrastructure():
            sys.exit(1)
        
        # Build and run the application
        if not build_and_run():
            sys.exit(1)
        
        print("\nMXShop application has been successfully deployed.")
        print("Web interfaces available at:")
        print("- User service: http://localhost:8021")
        print("- Goods service: http://localhost:8022")
        print("- Order service: http://localhost:8023")
        print("- User Operations service: http://localhost:8024")
        print("- OSS service: http://localhost:8025")
        print("\nNacos console: http://localhost:8848/nacos (username: nacos, password: nacos)")
        print("\nTo stop the application, run: docker-compose down")
    except KeyboardInterrupt:
        print("\nDeployment interrupted by user.")
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error during deployment: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
