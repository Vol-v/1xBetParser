package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"valentin-lvov/1x-parser/scrapper"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func StartConsumer(client *redis.Client, ch *amqp.Channel) error {

	q, err := ch.QueueDeclare("parse-1xbet-queue", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, "service.parse.1xbet", "service.parse", false, nil)

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var task ScrapingTask
			fmt.Printf("got message in queue!\n")
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("Error decoding message: %s", err)
				continue
			}

			log.Printf("Received a task: URL=%s, Duration=%d", task.URL, task.Duration)
			// Trigger scraping based on the task details
			go scrapper.TrackWebsite(task.URL, time.Second*time.Duration(task.Duration), time.Second*(10), client)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
