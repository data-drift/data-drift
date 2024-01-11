package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
)

type MetricStorageKey string

var redisClient *redis.Client

func getRedisClient() (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}
	var redisURL = os.Getenv("REDIS_TLS_URL")
	if redisURL == "" {
		redisURL = os.Getenv("REDIS_URL")
	}

	redisOpt, redisErr := redis.ParseURL(redisURL)
	redisOpt.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if redisErr != nil {
		return nil, redisErr
	}
	rdb := redis.NewClient(redisOpt)
	redisClient = rdb
	return rdb, nil
}

func ReadMetricKPI(path MetricStorageKey) (Metrics, error) {
	rdb, redisErr := getRedisClient()

	if redisErr != nil {
		jsonFile, err := os.ReadFile(string(path))
		if err != nil {
			return nil, err
		}

		var data Metrics
		err = json.Unmarshal(jsonFile, &data)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else {
		var ctx = context.Background() // TODO: use context from gin

		jsonData, err := rdb.Get(ctx, string(path)).Bytes()
		if err != nil {
			log.Printf("Could not get key. Err: %s", err)
			return nil, err
		}

		var data Metrics
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

func WriteMetricKPI(installationId int, metricName string, lineCountAndKPIByDateByVersion Metrics) MetricStorageKey {
	metricStoredFilePath := LegacyGetMetricStorageKey(fmt.Sprint(installationId), metricName)
	rdb, redisErr := getRedisClient()

	if redisErr != nil {

		file, err := os.Create(string(metricStoredFilePath))
		if err != nil {
			fmt.Printf("Error creating file: %v", err)
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		if err := enc.Encode(lineCountAndKPIByDateByVersion); err != nil {
			fmt.Printf("Error writing JSON to file: %v", err)
		}
		fmt.Println("Results written to file")
	} else {
		fmt.Println("Storing results in Redis")

		jsonData, err := json.Marshal(lineCountAndKPIByDateByVersion)
		if err != nil {
			fmt.Printf("Error occurred during marshaling. Err: %s", err)
		}
		var ctx = context.Background() // TODO: use context from gin

		err = rdb.Set(ctx, string(metricStoredFilePath), jsonData, 0).Err()
		if err != nil {
			fmt.Printf("Could not set key. Err: %s", err)
		}
	}
	return metricStoredFilePath
}

func LegacyGetMetricStorageKey(installationId string, metricName string) MetricStorageKey {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("dist/%s_%s_lineCountAndKPIByDateByVersion.json", installationId, metricNameEncoded)
	return MetricStorageKey(filepath)
}

func NewGetMetricStorageKey(owner, repo, metricName string) MetricStorageKey {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("dist/%s/%s/%s.json", owner, repo, metricNameEncoded)
	return MetricStorageKey(filepath)
}
