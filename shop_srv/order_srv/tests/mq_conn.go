package main

//func main() {
//	producer, err := rocketmq.NewProducer(
//		rocketmq.ProducerConfig{
//			NameServer: "127.0.0.1:9876",
//			GroupName:  "testGroup",
//		})
//	if err != nil {
//		fmt.Printf("create producer error: %s\n", err.Error())
//		return
//	}
//
//	err = producer.Start()
//	if err != nil {
//		fmt.Printf("start producer error: %s\n", err.Error())
//		return
//	}
//
//	defer producer.Shutdown()
//
//	transactionListener := &TransactionListener{producer: producer}
//
//	// 发送事务消息
//	err = producer.SendMessageInTransaction(
//		primitive.NewMessage("TopicTest", []byte("Hello RocketMQ")),
//		transactionListener,
//		"Hello RocketMQ")
//	if err != nil {
//		fmt.Printf("send message in transaction error: %s\n", err.Error())
//		return
//	}
//}
//
//// TransactionListener 事务监听器
//type TransactionListener struct {
//	producer *rocketmq.Producer
//}
//
//// ExecuteLocalTransaction 执行本地事务
//func (tl *TransactionListener) ExecuteLocalTransaction(msg *primitive.Message, arg interface{}) (rocketmq.TransactionStatus, error) {
//
//	// 执行本地事务
//	// ...
//
//	// 发送延时消息
//	delayProducer, err := rocketmq.NewProducer(
//		rocketmq.ProducerConfig{
//			NameServer: "127.0.0.1:9876",
//			GroupName:  "testGroup_delay",
//		})
//	if err != nil {
//		fmt.Printf("create delay producer error: %s\n", err.Error())
//		return rocketmq.TransactionStatusUnknown, err
//	}
//
//	err = delayProducer.Start()
//	if err != nil {
//		fmt.Printf("start delay producer error: %s\n", err.Error())
//		return rocketmq.TransactionStatusUnknown, err
//	}
//	defer delayProducer.Shutdown()
//
//	_, err = delayProducer.SendSync(primitive.NewMessage("TopicTest", []byte("Delay Message")))
//	if err != nil {
//		fmt.Printf("send delay message error: %s\n", err.Error())
//		return rocketmq.TransactionStatusRollback, err
//	}
//
//	return rocketmq.TransactionStatusCommit, nil
//}
