package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

// Init初始化
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer close fail,err:", err)
		return
	}

	logDataChan = make(chan *logData, maxSize)
	//开启后台的goroutine从通道中取数据发往kafka
	go sendToKafka()
	return
}

func SendToChan(topic, data string) {
	msg := new(logData)
	msg.topic = topic
	msg.data = data
	logDataChan <- msg
}

func sendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			msg := &sarama.ProducerMessage{}
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("send msg failde,err:", err)
			}
			fmt.Printf("pid:%v offset:%v/n", pid, offset)
		default:
			time.Sleep(time.Microsecond * 50)
		}
	}
}
