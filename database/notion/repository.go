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

const DATADRIFT_PROPERTY = "datadrift-id"

func FindOrCreateReportPageId(apiKey string, databaseId string, reportName string) (string, error) {
	err := AssertDatabaseHasDatadriftProperties(databaseId, apiKey)
	if err != nil {
		return "", err
	}
	existingReportId, err := QueryDatabaseWithReportId(apiKey, databaseId, reportName)
	if err != nil {
		return "", err
	}
	if existingReportId == "" {
		fmt.Println("No existing report found, creating new one")
		return "New report id", nil
	}
	return existingReportId, nil
}

type httpTransport struct {
	w io.Writer
}

func QueryDatabaseWithReportId(apiKey string, databaseId string, reportId string) (string, error) {
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	queryParams := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "datadrift-id",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				RichText: &notion.TextPropertyFilter{
					Equals: reportId,
				},
			},
		},
	}

	existingReport, err := client.QueryDatabase(ctx, databaseId, queryParams)
	if err != nil {
		return "", err
	}
	fmt.Println("Number of existing report", len(existingReport.Results))
	if len(existingReport.Results) == 0 {
		return "", nil
	}
	if len(existingReport.Results) > 1 {
		fmt.Println("Warning: too many report with same id, returning first one")
	}
	return existingReport.Results[0].ID, nil
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
