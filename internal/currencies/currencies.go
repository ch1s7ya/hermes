package currencies

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type ExRate struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int                `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

// TODO: Neet to check reponse code 403
func GetAPIExRate() (*ExRate, error) {
	token := os.Getenv("OPENEXCHANGERATE_TOKEN")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://openexchangerates.org/api/latest.json", nil)

	if err != nil {
		log.Println("Request error:", err)
	}

	req.Header.Add("Authorization", "Token "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Response error:", err)
	}
	exchangeRate, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Can't get exchange rate from response body:", err)
	}

	var data ExRate

	err = json.Unmarshal(exchangeRate, &data)

	if err != nil {
		log.Println("Can't unmarshal from response body:", err)
	}

	return &data, nil
}

func CreateRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Can't connect to redis: %s", err)
	}

	return client
}

func GetCachedExRate(redisClient *redis.Client) (*ExRate, error) {
	ctx := context.Background()
	cachedData, err := redisClient.Get(ctx, "exchange_rate").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var data ExRate
	err = json.Unmarshal([]byte(cachedData), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func CacheExRate(redisClient *redis.Client, data *ExRate, ttl time.Duration) error {
	dataJSON, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return redisClient.Set(context.Background(), "exchange_rate", dataJSON, ttl).Err()
}

func GetExRate() *ExRate {
	client := CreateRedisClient()

	exrate, err := GetCachedExRate(client)
	if err != nil {
		log.Println("Some error", err)
	}

	if exrate == nil {
		exrate, err = GetAPIExRate()
		if err != nil {
			log.Println("Can't get exchange rate from API", err)
		}

		err := CacheExRate(client, exrate, time.Hour*24)
		if err != nil {
			log.Println("Can't cache exchange rate", err)
		}
		log.Println("Save exchange rate to cache")
	}
	return exrate
}
