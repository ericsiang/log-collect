## Go 併發程式設計練手項目 - 日志收集系统
參考：
* [日誌收集專案架構設計](https://github.com/moxi624/LearningNotes/tree/master/Golang/Golang%E8%BF%9B%E9%98%B6/17_%E6%97%A5%E5%BF%97%E6%94%B6%E9%9B%86%E9%A1%B9%E7%9B%AE%E6%9E%B6%E6%9E%84%E8%AE%BE%E8%AE%A1%E5%8F%8AKafka%E4%BB%8B%E7%BB%8D)

* [日志收集系统](https://blog.csdn.net/qq_73924465/category_12643285.html)

* [Go实现LogCollect](https://blog.csdn.net/weixin_45565886/article/details/132630758?ops_request_misc=%257B%2522request%255Fid%2522%253A%252213af6ecff21c4730634ce6fd2bd36658%2522%252C%2522scm%2522%253A%252220140713.130102334..%2522%257D&request_id=13af6ecff21c4730634ce6fd2bd36658&biz_id=0&utm_medium=distribute.pc_search_result.none-task-blog-2~all~sobaiduend~default-4-132630758-null-null.142^v101^pc_search_result_base8&utm_term=go%20%E6%97%A5%E5%BF%97%E6%94%B6%E9%9B%86%E7%B3%BB%E7%BB%9F&spm=1018.2226.3001.4187)

### etcd
參考：
* [Docker Compose 部署 etcd 集群](https://oldme.net/article/32)
* [golang etcd 简明教程](https://learnku.com/articles/37343)

docker-compose.yml
```
services:
  etcd1:
    image: "bitnami/etcd:latest"
    container_name: "etcd1"
    # restart: "always"
    ports:
      - 12379:12379
    environment:
      - TZ=Asia/Taipei
      # 允许无认证访问
      - ALLOW_NONE_AUTHENTICATION=yes
      # 名称
      - ETCD_NAME=etcd1
      # 列出这个成员的伙伴 URL 以便通告给集群的其他成员
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd1:2380
      # 用于监听伙伴通讯的URL列表
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      # 用于监听客户端通讯的URL列表
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:12379
      # 初始化集群记号，建议一个集群内的所有节点都设置唯一的token
      - ETCD_INITIAL_CLUSTER_TOKEN=MyEtcd
      # 列出这个成员的客户端URL，通告给集群中的其他成员
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd1:12379
      # 集群配置，格式是ETCD_NAME=ETCD_INITIAL_ADVERTISE_PEER_URLS
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      # 初始化集群状态，这里写 new 就行
      - ETCD_INITIAL_CLUSTER_STATE=new

  etcd2:
    image: "bitnami/etcd:latest"
    container_name: "etcd2"
    # restart: "always"
    ports:
      - 22379:22379
    environment:
      - TZ=Asia/Taipei
      # 允许无认证访问
      - ALLOW_NONE_AUTHENTICATION=yes # 名称
      - ETCD_NAME=etcd2
      # 列出这个成员的伙伴 URL 以便通告给集群的其他成员
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd2:2380
      # 用于监听伙伴通讯的URL列表
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      # 用于监听客户端通讯的URL列表
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:22379
      # 初始化集群记号，建议一个集群内的所有节点都设置唯一的token
      - ETCD_INITIAL_CLUSTER_TOKEN=MyEtcd
      # 列出这个成员的客户端URL，通告给集群中的其他成员
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd2:22379
      # 集群配置，格式是ETCD_NAME=ETCD_INITIAL_ADVERTISE_PEER_URLS
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      # 初始化集群状态，这里写 new 就行
      - ETCD_INITIAL_CLUSTER_STATE=new

  etcd3:
    image: "bitnami/etcd:latest"
    container_name: "etcd3"
    # restart: "always"
    ports:
      - 32379:32379
    environment:
      - TZ=Asia/Taipei
      # 允许无认证访问
      - ALLOW_NONE_AUTHENTICATION=yes # 名称
      - ETCD_NAME=etcd3
      # 列出这个成员的伙伴 URL 以便通告给集群的其他成员
      - ETCD_INITIAL_ADVERTISE_PEER_URLS=http://etcd3:2380
      # 用于监听伙伴通讯的URL列表
      - ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380
      # 用于监听客户端通讯的URL列表
      - ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:32379
      # 初始化集群记号，建议一个集群内的所有节点都设置唯一的token
      - ETCD_INITIAL_CLUSTER_TOKEN=MyEtcd
      # 列出这个成员的客户端URL，通告给集群中的其他成员
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd3:32379
      # 集群配置，格式是ETCD_NAME=ETCD_INITIAL_ADVERTISE_PEER_URLS
      - ETCD_INITIAL_CLUSTER=etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380
      # 初始化集群状态，这里写 new 就行
      - ETCD_INITIAL_CLUSTER_STATE=new
  etcdkeeper:
    image: evildecay/etcdkeeper
    ports:
      - 8088:8080
```
進到 etcd1 容器
```
docker exec -it etcd1 bash
```
容器內下 etcd 指令
```
etcdctl --endpoints=localhost:12379 put  /t v
etcdctl --endpoints=etcd2:22379 get /t
etcdctl --endpoints=etcd3:32379 get /t
```

### kafka
docker-compose.yml
```
services:
  zookeeper:
    image: 'bitnami/zookeeper:latest'
    ports:
      - '2181:2181'
    environment:
      # 匿名登录--必须开启
      - ALLOW_ANONYMOUS_LOGIN=yes
      #volumes:
      #- ./zookeeper:/bitnami/zookeeper
      # 该镜像具体配置参考 https://github.com/bitnami/bitnami-docker-kafka/blob/master/README.md
  kafka:
    image: 'bitnami/kafka:latest'
    ports:
      - '9092:9092'
      - '9999:9999'
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_CFG_NODE_ID=1
      # - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092
      # - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      # - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      # - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT
      # - KAFKA_CFG_ZOOKEEPER_METADATA_MIGRATION_ENABLE=true
      # 客户端访问地址，更换成自己的
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      # 允许使用PLAINTEXT协议(镜像中默认为关闭,需要手动开启)
      - ALLOW_PLAINTEXT_LISTENER=yes
      # 开启自动创建 topic 功能便于测试
      - KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true
      # 全局消息过期时间 6 小时(测试时可以设置短一点)
      - KAFKA_CFG_LOG_RETENTION_HOURS=6
      # 开启JMX监控
      # - JMX_PORT=9999
      #volumes:
      #- ./kafka:/bitnami/kafka
    # depends_on:
      # - zookeeper
  # Web 管理界面 另外也可以用exporter+prometheus+grafana的方式来监控 https://github.com/danielqsj/kafka_exporter
  # kafka_manager 只有支援 linux/amd64
  # kafka_manager:
  #   image: 'hlebalbau/kafka-manager:latest'
  #   ports:
  #     - "9000:9000"
  #   environment:
  #     ZK_HOSTS: "zookeeper:2181"
  #     APPLICATION_SECRET: letmein
  #   depends_on:
  #     - zookeeper
  #     - kafka
  # kafka-ui:
  #   image: provectuslabs/kafka-ui
  #   container_name: kafka-ui
  #   ports:
  #     - "8082:8080"
  #   environment:
  #     - KAFKA_CLUSTERS_0_NAME=local-kafka
  #     - KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092
  #     - KAFKA_CLUSTERS_0_ZOOKEEPER=zookeeper:2181
  #     # - KAFKA_CLUSTERS_0_READONLY=true
  #   depends_on:
  #     - kafka
  #     - zookeeper
  Redpanda:
    image: redpandadata/console:latest
    ports:
      - "8083:8080"
    environment:
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - kafka
      - zookeeper

```


### elk 
docker-compose.yml
```
services:
  elasticsearch:
    image: elasticsearch:8.17.1
    container_name: es
    ports:
      - 9200:9200
      - 9300:9300
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=true
      - ELASTIC_PASSWORD=123456
      - ES_JAVA_OPTS=-Xms128m -Xmx128m # 出現 error code exit(137) 加此行
         # 開發測試用，關閉 SSL
      - xpack.security.http.ssl.enabled=false
      - xpack.security.transport.ssl.enabled=false
      # - "logger.level=DEBUG"
    networks:
      - elastic
  kibana:
    image: kibana:8.17.1
    container_name: kibana
    ports:
      - "5601:5601"
    # volumes:
    #   - ./kibana.yml:/usr/share/kibana/config/kibana.yml
    environment:
      - ELASTICSEARCH_HOSTS="http://elasticsearch:9200"
      # - XPACK_SECURITY_ENABLED=false
      # 創建一個專門的 Kibana 系統用戶
      - ELASTICSEARCH_USERNAME=kibana_system
      - ELASTICSEARCH_PASSWORD=123456
    depends_on:
      - elasticsearch
    networks:
      - elastic
networks:
  elastic:
    driver: bridge
```