package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	proto "order_srv/proto/v1order"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.1.106:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)
}

func main() {
	Init()
	//TestCreateCartItem(12, 1, 422)
	//TestCartItemList(12)
	//TestUpdateCartItem(1)
	//TestCreateOrder()
	//TestGetOrderDetail(1)
	TestOrderList()
	conn.Close()
}
