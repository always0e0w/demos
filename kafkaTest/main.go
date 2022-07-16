package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func main() {
	// msg := fmt.Sprintf("test message produced on %s", time.Now().Format("2006-01-02 15:04:05"))
	// producerMessage(msg)
	// consumeMessage()
	consumeMessageWithGroup()
}

func producerMessage(message string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	msg := sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(&msg)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(partition, offset)
	}
}

func consumeMessage() {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	partitionConsumer, err := consumer.ConsumePartition("test_topic", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	consumed := 0
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("message: %s, offset:%d\n", msg.Value, msg.Offset)
			consumed += 1
		case <-signals:
			log.Printf("consumed %d messages. return now.\n", consumed)
			return
		}
	}
}

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func consumeMessageWithGroup() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "test_consumer_group", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("[ERROR]", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{"test_topic"}
		handler := exampleConsumerGroupHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
