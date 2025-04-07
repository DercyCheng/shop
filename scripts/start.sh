
#!/bin/sh

# Start microservices in background
echo "Starting MXShop microservices..."

# Start service layer first
/app/bin/user_srv &
/app/bin/goods_srv &
/app/bin/inventory_srv &
/app/bin/order_srv &
/app/bin/userop_srv &

# Wait for services to be ready
sleep 10

# Start web layer
/app/bin/user_web &
/app/bin/goods_web &
/app/bin/order_web &
/app/bin/userop_web &
/app/bin/oss_web &

# Keep container running
echo "All services started. Container is now running."
tail -f /dev/null