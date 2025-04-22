package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}
func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Date {
		fmt.Println(user.Mobile, user.NickName, user.Password)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "password",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}
func TestCresteUser(name, passwd, mobile string) {
	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: name,
		PassWord: passwd,
		Mobile:   mobile,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)

}
func main() {
	Init()
	TestGetUserList()
	conn.Close()
	//for i := 1; i < 10; i++ {
	//	TestCresteUser("jzin"+strconv.Itoa(i), "password", "1532509548"+strconv.Itoa(i))
	//}
}
