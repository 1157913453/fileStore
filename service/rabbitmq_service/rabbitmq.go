package rabbitmq_service

import (
	cfg "filestore/config"
	"filestore/service/oss_service"
	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	MQConn *amqp.Connection
	Ch     *amqp.Channel
)

func init() {
	var err error
	MQConn, err = amqp.Dial(cfg.RabbitURL)
	if err != nil {
		log.Errorf("初始化RabbitMQ失败：%v", err)
	} else {
		log.Infof("初始化RabbitMQ成功")
	}
	Ch, err = MQConn.Channel()
	if err != nil {
		log.Errorf("创建RabbitMQ通道失败:%v", err)
	}
	go ReceiveMQ()
}

type SendMQMsg struct {
	FileAddr string
	ChunkNum int
}

func SendMQ(fileAddr string, chunkNum int) error {
	msg := SendMQMsg{fileAddr, chunkNum}
	body, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("序列化msg失败：%v", err)
		return err
	}
	err = Ch.Publish(cfg.TransExchangeName, cfg.TransOSSRoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		log.Errorf("向RabbitMQ中发送信息失败：%v", err)
	}
	return err
}

func ReceiveMQ() {
	msgs, err := Ch.Consume(cfg.TransOSSQueueName, "transfer_oss", true, false, false, false, nil)
	if err != nil {
		log.Errorf("接受RabbitMQ消息错误：%v", err)
		return
	}

	for msg := range msgs {
		// 起个goroutine执行任务
		m := &SendMQMsg{}
		err = json.Unmarshal(msg.Body, m)
		if err != nil {
			log.Errorf("反序列化失败：%v", err)
			return
		}
		go func() {
			err = oss_service.OssUploadPart(m.FileAddr, m.ChunkNum)
			if err != nil {
				log.Errorf("上传OSS失败：%v", err)
				return
			}
		}()

	}

}
