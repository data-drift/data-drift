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

func GetRedisClient() (*redis.Client, error) {

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
	return rdb, nil
}

type KpiRepository struct {
	RedisClient *redis.Client
}

func NewKpiRepository(redisClient *redis.Client) *KpiRepository {
	return &KpiRepository{RedisClient: redisClient}
}

func (h *KpiRepository) ReadMetricKPI(path MetricStorageKey) (Metrics, error) {
	var ctx = context.Background() // TODO: use context from gin

	jsonData, err := h.RedisClient.Get(ctx, string(path)).Bytes()
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

func (h *KpiRepository) WriteMetricKPI(repoOwner string, repoName string, metricName string, lineCountAndKPIByDateByVersion Metrics) MetricStorageKey {
	metricStoredFilePath := NewGetMetricStorageKey(repoOwner, repoName, metricName)

	jsonData, err := json.Marshal(lineCountAndKPIByDateByVersion)
	if err != nil {
		fmt.Printf("Error occurred during marshaling. Err: %s", err)
	}
	var ctx = context.Background() // TODO: use context from gin

	err = h.RedisClient.Set(ctx, string(metricStoredFilePath), jsonData, 0).Err()
	if err != nil {
		fmt.Printf("Could not set key. Err: %s", err)
	}
	return metricStoredFilePath
}

func LegacyGetMetricStorageKey(installationId string, metricName string) MetricStorageKey {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("dist/%s_%s_lineCountAndKPIByDateByVersion.json", installationId, metricNameEncoded)
	log.Println("Using legacy storage key" + filepath)
	return MetricStorageKey(filepath)
}

func NewGetMetricStorageKey(owner, repo, metricName string) MetricStorageKey {
	metricNameEncoded := url.PathEscape(metricName)
	filepath := fmt.Sprintf("%s/%s/%s", owner, repo, metricNameEncoded)
	return MetricStorageKey(filepath)
}
