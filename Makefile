SHELL=cmd.exe

zookeeper:
	D:\kafka\kafka_2.13-3.4.0\bin\windows\zookeeper-server-start.bat D:\kafka\kafka_2.13-3.4.0\config\zookeeper.properties

kafka:
	D:\kafka\kafka_2.13-3.4.0\bin\windows\kafka-server-start.bat D:\kafka\kafka_2.13-3.4.0\config\server.properties

deleteKafkaLogs:
	rd /Q /S D:\kafka\kafkalogs

server:
	go run main.go

.PHONY:zookeeper kafka redis server