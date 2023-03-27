SHELL=cmd.exe

zookeeper:
	D:\Kafka\kafka_2.13-3.4.0\bin\windows\zookeeper-server-start.bat D:\Kafka\kafka_2.13-3.4.0\config\zookeeper.properties

kafka:
	D:\Kafka\kafka_2.13-3.4.0\bin\windows\kafka-server-start.bat D:\Kafka\kafka_2.13-3.4.0\config\server.properties

deleteKafkaLogs:
	rd /Q /S D:\kafka\kafkalogs

server:
	go run main.go

.PHONY:zookeeper kafka redis server