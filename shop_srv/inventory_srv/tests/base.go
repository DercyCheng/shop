package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
)

var brandClient proto.InventoryClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.1.106:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()
	//var i int32
	//for i = 421; i <= 840; i++ {
	//	TestSetInv(i, 100)
	//}

	//TestInvDetail(421)

	//
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
	//TestReback()
	conn.Close()
}
