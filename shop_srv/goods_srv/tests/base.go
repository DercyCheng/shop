package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Init() {
	var err error
	conn, err = grpc.Dial("10.231.72.37:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	brandClient = proto.NewGoodsClient(conn)
}
func main() {
	Init()
	//TestGetCategoryList()
	//TestGetSubCategoryList()
	//TestGetCategoryBrandList()
	//TestGetGoodsList()
	//TestGetBatchGoods()
	TestGetGoodsDetail()
	conn.Close()
	//for i := 1; i < 10; i++ {
	//	TestCresteUser("jzin"+strconv.Itoa(i), "password", "1532509548"+strconv.Itoa(i))
	//}
}
