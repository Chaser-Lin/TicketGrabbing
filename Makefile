SHELL=cmd.exe

zookeeper:
	D:\Kafka\kafka_2.13-3.4.0\bin\windows\zookeeper-server-start.bat D:\Kafka\kafka_2.13-3.4.0\config\zookeeper.properties

kafka:
	D:\Kafka\kafka_2.13-3.4.0\bin\windows\kafka-server-start.bat D:\Kafka\kafka_2.13-3.4.0\config\server.properties

redis:
	D:\Redis\Redis-x64-5.0.10\redis-server.exe --service-start

server:
	go run main.go

.PHONY:zookeeper kafka redis server