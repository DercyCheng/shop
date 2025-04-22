package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}
func TestGetCategoryBrandList() {
	rsp, err := brandClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 135475,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
	//for _, brand := range rsp.Data {
	//	fmt.Println(brand.Name)
	//}
}
