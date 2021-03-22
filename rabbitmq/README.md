## RabbitMQ

Step 1:<br/> 
Create a local docker network *rabbits* to enable instances
created to communicate.<br/> 
In addition, both the publisher and consumer
apps would also run on this network.<br/>

```bash
$ docker network create rabbits
```

Step 2:<br/>
Run a container on the network in background mode specifying the 
hostname and official RabbitMQ image.<br/>
```bash
$ docker run -d --rm --net rabbits -p 8080:15672 --hostname rabbit-1 --name rabbit-1 rabbitmq:3.8-management
```

Step 3:<br/>
Check docker logs to see if your RabbitMQ instance started up successfully 
```bash
$ docker logs rabbit-1
```

Step 4:<br/>
Use the RabbitMQ command line tool *rabbitmqctl* to manage the node.<br/>
An additional command line tool to RabbitMQ, *rabbitmq-plugins* is used integrate and manage plugins.<br/>
```bash
$ docker exec -it rabbit-1 bash    
$ rabbitmqctl [--node <node>] [--timeout <timeout>] [--longnames] [--quiet] <command> [<command options>]
$ rabbitmq-plugins [--node <node>] [--timeout <timeout>] [--longnames] [--quiet] <command> [<command options>]
```

Step 5:<br/>
Build and run the publisher application.<br/>
```bash
$ cd messaging/rabbitmq/app/publisher
$ docker build . -t fokorji/rabbitmq-publisher:v1.0.0
$ docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest -p 80:80 fokorji/rabbitmq-publisher:v1.0.0
```

Step 6:<br/>
Build and run the consumer application.<br/>
```bash
$ cd messaging/rabbitmq/app/consumer
$ docker build . -t fokorji/rabbitmq-consumer:v1.0.0
$ docker run -it --rm --net rabbits -e RABBIT_HOST=rabbit-1 -e RABBIT_PORT=5672 -e RABBIT_USERNAME=guest -e RABBIT_PASSWORD=guest fokorji/rabbitmq-consumer:v1.0.0
```
