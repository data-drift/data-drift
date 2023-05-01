package notion_database

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dstotijn/go-notion"
)

const DATADRIFT_PROPERTY = "datadrift-id4"

func FindOrCreateReportPageId(apiKey string, databaseId string, reportName string) (string, error) {
	err := AssertDatabaseHasDatadriftProperties(databaseId, apiKey)
	if err != nil {
		return "", err
	}
	return databaseId, nil
}

type httpTransport struct {
	w io.Writer
}

func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}

func AssertDatabaseHasDatadriftProperties(databaseID, apiKey string) error {
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))
	database, err := client.FindDatabaseByID(ctx, databaseID)

	hasDatadriftProperty := false

	for _, property := range database.Properties {
		fmt.Println("Property:", property.Name)
		if property.Name == DATADRIFT_PROPERTY {
			hasDatadriftProperty = true
		}

	}
	fmt.Println("hasDatadriftProperty:", hasDatadriftProperty)
	if !hasDatadriftProperty {
		params := notion.UpdateDatabaseParams{
			Properties: map[string]*notion.DatabaseProperty{
				DATADRIFT_PROPERTY: {
					Type:     "rich_text",
					RichText: &notion.EmptyMetadata{},
				},
			},
		}

		fmt.Println("Creating property", params)
		updatedDB, err := client.UpdateDatabase(ctx, databaseID, params)
		if err != nil {
			return err
		}
		fmt.Println("Updated database", updatedDB, " with property", DATADRIFT_PROPERTY)
	}
	return err
}
