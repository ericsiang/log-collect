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