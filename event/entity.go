package event

type Config struct {
	Address string `yaml:"Address"`
	Topic   string `yaml:"Topic"`
}

type Message struct {
	TicketID int
	UserID   int
}

var kafkaConf Config

func Init(config Config) {
	kafkaConf = config
}
