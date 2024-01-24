package queue

import (
	"context"
	"encoding/json"
	"time"
	"valentin-lvov/1x-parser/cache"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type ScrapingTask struct {
	URL      string `json:"url"`
	Duration int    `json:"duration"`
}

func ConnectToRabbitMQ(uri string) (*amqp.Connection, *amqp.Channel, error) {
	//TODO: make connection to uri from config
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	err = ch.ExchangeDeclare(
		"service.parse", // Exchange name
		"direct",        // Type of the exchange
		true,            // Durable
		false,           // Auto-deleted
		false,           // Internal
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}

func PublishTrackingTask(url string, duration int, ch *amqp.Channel) error {
	// TODO: make connection to uri from config

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	task := ScrapingTask{
		URL:      url,
		Duration: duration,
	}

	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, "service.parse", "service.parse.1xbet", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	})
	return err
}

func GetTrackingResults(rdb *redis.Client, url string) (map[string]string, error) {
	/*separate function in case I need to Get tracking results from different places.
	Then I'll just implement and inject the interface "resultGetter" or something*/

	/*get results from redis*/
	result, err := cache.RetrieveFromRedis(rdb, url)
	return result, err
}
