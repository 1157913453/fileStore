package rabbitmq_service

import (
	cfg "filestore/config"
	"filestore/service/oss_service"
	json "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
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
	}
	Ch, err = MQConn.Channel()
	if err != nil {
		log.Errorf("创建RabbitMQ通道失败:%v", err)
	}

	err = Ch.ExchangeDeclare(cfg.ExchangeName, cfg.ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Errorf("创建RabbitMQ交换机错误")
		panic(err)
	}

	_, err = Ch.QueueDeclare(cfg.QueueName, cfg.QueueDurable, false, false, false, nil)
	if err != nil {
		log.Errorf("创建RabbitMQ队列错误")
		panic(err)
	}

	err = Ch.QueueBind(cfg.QueueName, cfg.RoutingKey, cfg.ExchangeName, false, nil)
	if err != nil {
		log.Errorf("交换机和Queue绑定错误")
		panic(err)
	}
	log.Infof("初始化RabbitMQ成功")
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
	err = Ch.Publish(cfg.ExchangeName, cfg.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		log.Errorf("向RabbitMQ中发送信息失败：%v", err)
	}
	return err
}

func ReceiveMQ() {
	msgs, err := Ch.Consume(cfg.QueueName, "transfer_oss", true, false, false, false, nil)
	if err != nil {
		log.Errorf("接受RabbitMQ消息错误：%v", err)
		return
	}
	limitChan := make(chan struct{}, 1000) // 最多同时存在1000个上传oss的任务
	for msg := range msgs {
		// 起个goroutine执行任务
		m := &SendMQMsg{}
		err = json.Unmarshal(msg.Body, m)
		if err != nil {
			log.Errorf("反序列化失败：%v", err)
			return
		}
		limitChan <- struct{}{}
		go func() {
			err = oss_service.OssUploadPart(m.FileAddr, m.ChunkNum)
			if err != nil {
				log.Errorf("上传OSS失败：%v", err)
				return
			}
			<-limitChan
			log.Infof("文件%s上传OSS成功", m.FileAddr)
			// 上传OSS成功后删除本地文件
			err = os.Remove(m.FileAddr)
			if err != nil {
				log.Errorf("删除本地文件失败：%v", err)
				return
			}
		}()

	}

}
