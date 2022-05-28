package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func sub(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic %s\n", topic)
}

func publish(client mqtt.Client, message string) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func main() {
	var broker string
	var port = 1883
	fmt.Println("请输入EMQX服务器的IP地址")

	fmt.Scan(&broker)

	fmt.Printf("Broker IP: tcp://%s:%d\n", broker, port)
	name, _ := os.Hostname()
	fmt.Printf("ClientID: %s\n", name)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(name + "-c")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	// sub(client, topic)
	// fmt.Printf("已连接。\n")
	time.Sleep(1)
	fmt.Printf("输入 quit 并回车可退出程序\n")
	var enter string
	for {
		fmt.Println("请指定topic")
		var topic string
		fmt.Scan(&topic)
		fmt.Println("请输入要发送的消息")
		fmt.Scan(&enter)
		token := client.Publish(topic, 0, false, enter)
		token.Wait()
		fmt.Println("发送成功")
		if enter == "quit" || topic == "quit" {
			return
		}
	}

}
