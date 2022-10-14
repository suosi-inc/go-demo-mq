# go-demo-mq

## RabbitMQ 配置说明

```
# 服务端行为
rabbit:
  amqp: amqp://admin:admin@localhost:5672  
  prefetch: 1  # 消息预取数量
  deliveryMode: 2  # 消息投递持久化, 持久化2, 非持久化1, 默认是 2
    
# 客户端行为
rabbitQueue:
  simple:
    name: go-simple  # 队列名称
    goroutines: 3  # 协程数量
  topic:
    exchange: go-topic  # 交换机名称
    name: go-topic-log  # 队列名称
    goroutines: 3  # 协程数量
    routingKeys:  # routingKeys 切片
      - '#'  
```

## Kafka 配置说明

```
# 服务端行为
kafka:
  servers:
    - kafka-online1:9092
    - kafka-online2:9092
    - kafka-online3:9092
  compress: lz4  # 生产者压缩算法(gzip,snappy,zstd)，默认 lz4
  write:
    async: true  # 生产者消息确认方式，默认异步
  read:
    offset: last    # 新的消费者组读取消息的偏移模式(first,last), 默认 last
    commitInterval: 1   # ack 异步确认间隔秒数，0 表示同步确认, 默认是 0 
    MaxBytes: 4096000 # 消息最大字节数限制，应当小于服务器配置
  sasl:  # Sasl 认证相关配置
    enable: false
    mechanism: PLAIN
    username: admin
    password: admin

# 客户端行为
kafkaQueue:
  topic: go-test # Topic
  groupId: go-test-group  # 消费者组
  goroutines: 3  # # 协程数量

```