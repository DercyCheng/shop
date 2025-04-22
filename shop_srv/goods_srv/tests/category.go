package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestGetCategoryList() {
	rsp, err := brandClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	//var ext_de_json map[string]interface{}
	//json.Unmarshal([]byte(rsp.JsonData), &ext_de_json)
	fmt.Println(rsp.JsonData)
	//for _, category := range rsp.Data {
	//	fmt.Println(category.Name)
	//}
}
func TestGetSubCategoryList() {
	rsp, err := brandClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 130358,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.SubCategorys)
}
