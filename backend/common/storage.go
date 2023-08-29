package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
)

type MetricRedisKey string

var ctx = context.Background()

func ReadMetricKPI(path MetricRedisKey) (Metrics, error) {
	var redisURL = os.Getenv("REDIS_URL")
	redisOpt, redisErr := redis.ParseURL(redisURL)

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
		var rdb = redis.NewClient(redisOpt)

		jsonData, err := rdb.Get(ctx, string(path)).Bytes()
		if err != nil {
			log.Fatalf("Could not get key. Err: %s", err)
			return nil, err
		}

		var data Metrics
		err = json.Unmarshal(jsonData, &data)
		if err != nil {
			return nil, err
		}
		rdb.Close()
		return data, nil
	}
}

func WriteMetricKPI(installationId int, metricName string, lineCountAndKPIByDateByVersion Metrics) MetricRedisKey {
	metricStoredFilePath := GetMetricStorageKey(fmt.Sprint(installationId), metricName)
	var redisURL = os.Getenv("REDIS_URL")
	redisOpt, redisErr := redis.ParseURL(redisURL)

	if redisErr != nil {

		file, err := os.Create(string(metricStoredFilePath))
		if err != nil {
			log.Fatalf("Error creating file: %v", err.Error())
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		if err := enc.Encode(lineCountAndKPIByDateByVersion); err != nil {
			log.Fatalf("Error writing JSON to file: %v", err.Error())
		}
		fmt.Println("Results written to file")
	} else {
		var rdb = redis.NewClient(redisOpt)
		fmt.Println("Storing results in Redis")

		jsonData, err := json.Marshal(lineCountAndKPIByDateByVersion)
		if err != nil {
			log.Fatalf("Error occurred during marshaling. Err: %s", err)
		}
		err = rdb.Set(ctx, string(metricStoredFilePath), jsonData, 0).Err()
		if err != nil {
			log.Fatalf("Could not set key. Err: %s", err)
		}
		rdb.Close()
	}
	return metricStoredFilePath
}

func GetMetricStorageKey(installationId string, metricName string) MetricRedisKey {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("dist/%s_%s_lineCountAndKPIByDateByVersion.json", installationId, metricNameEncoded)
	return MetricRedisKey(filepath)
}
