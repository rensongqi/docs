# 1 RocketMQ架构介绍

![image-20210719135055312](C:\Users\songqi\AppData\Roaming\Typora\typora-user-images\image-20210719135055312.png)

具体操作流程如下：

1、启动Namesrv，Namesrv起来后监听端口，等待Broker、Producer、Consumer连上来，相当于一个路由控制中心
2、Broker启动，跟所有的Namesrv保持长连接，定时发送心跳包。心跳包中包含当前Broker信息（IP+端口等）以及存储所有Topic信息
3、收发消息前，先创建topic，创建topic时需要指定该topic要存储在哪些Broker上，也可以在发送消息时自动创建Topic。
4、Producer发送消息，启动时先跟Namesrv集群中的其中一台建立长连接，并从Namesrv中获取当前发送的Topic存在哪些Broker上
5、Consumer跟Producer类似。

![image-20210719135121614](C:\Users\songqi\AppData\Roaming\Typora\typora-user-images\image-20210719135121614.png)

使用RabbitMQ可以发送普通消息、顺序消息、事务消息，顺序消息能实现有序消费，事务消息可以解决分布式事务实现数据一致。

RocketMQ有两种常见的消费模式，分别是DefaultMQPushConsumer和DefaultMQPullConsumer模式，这两种模式字面理解是一个是推送消息，一个是拉取消息。这里有个误区，其实无论是Push还是Pull，其本质都是拉取消息，只是实现机制不一样。

DefaultMQPushConsumer其实并不是Broker主动向Consumer推送消息，而是Consumer向Broker发送请求，保持一种长连接，Broker会每5s检测一次是否有消息，如果有消息，则将消息推送给Consumer。使用DefaultMQPushConsumer实现消息消费，Broker会主动记录消息消费的偏移量。

DefaultMQPullConsumer是消费方主动去Broker拉取数据，一般会在本地使用定时任务实现，使用它获得消息状态方便、负载均衡性能可控，但是消息的及时性差，而且需要手动记录消息消费的偏移量信息，所以在工作中多数情况推荐使用Push模式。

RocketMQ发送的消息默认会存储到4个队列中，当然创建几个队列存储数据，可以自己定义。

消息发送有几个步骤：

1. 创建DefaultMQProducer
2. 设置Namesrv地址
3. 开启DefaultMQProducer
4. 创建消息Message
5. 发送消息
6. 关闭DefaultMQProducer