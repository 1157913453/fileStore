package config

const (
	// AsyncTransferEnable 是否开启文件异步转移（默认同步）
	AsyncTransferEnable = true
	// RabbitURL rabbitmq服务的url入口
	RabbitURL = "amqp://guest:guest@localhost:5672/"
	// TransExchangeName 交换机名字
	TransExchangeName = "uploadserver.trans"
	// TransOSSQueueName oss转移队列名
	TransOSSQueueName = "uploadserver.trans.oss"
	// TransOSSErrQueueName oss转移失败转入另一个队列名字
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	// TransOSSRoutingKey routingkey
	TransOSSRoutingKey = "oss"
)
