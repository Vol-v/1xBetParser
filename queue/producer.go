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

func PublishTrackingTask(url string, duration int) error {
	// TODO: make connection to uri from config
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

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
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
