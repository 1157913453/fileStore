package config

const (
	// AsyncTransferEnable 是否开启文件异步转移（默认同步）
	// AsyncTransferEnable = true
	// RabbitURL rabbitmq服务的url入口
	RabbitURL = "amqp://guest:guest@localhost:5672/"
	// ExchangeName 交换机名字
	ExchangeName = "uploadserver.trans"
	// ExchangeType 交换机类型（广播，单薄，组播，Header）
	ExchangeType = "direct"
	// TransOSSErrQueueName oss转移失败转入另一个队列名字
	TransOSSErrQueueName = "uploadserver.trans.oss.err"
	// RoutingKey routingkey
	RoutingKey = "oss"
	// QueueName oss转移队列名
	QueueName = "uploadserver.trans.oss"
	// QueueDurable 队列是否持久化
	QueueDurable = true
)
