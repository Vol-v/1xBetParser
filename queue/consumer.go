package queue

import (
	"encoding/json"
	"log"
	"valentin-lvov/1x-parser/cache"
	"valentin-lvov/1x-parser/scrapper"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	// scrap "valentin-lvov/1x-parser/scrapper"
)

func StartConsumer(client *redis.Client) error {
	// ... connection and channel setup ...

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("scrapeQueue", false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var task ScrapingTask
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("Error decoding message: %s", err)
				continue
			}

			log.Printf("Received a task: URL=%s, Duration=%d", task.URL, task.Duration)
			// Trigger scraping based on the task details
			data, err := scrapper.ScrapWebsite(task.URL)
			cache.StoreInRedis(client, task.URL, data)

		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
