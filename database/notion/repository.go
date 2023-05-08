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

var DefaultPropertiesToDelete = []string{"Tags", "Status", "Étiquette", "Étiquettes"}

func FindOrCreateReportPageId(apiKey string, databaseId string, reportName string) (string, error) {
	existingReportId, err := QueryDatabaseWithReportId(apiKey, databaseId, reportName)
	if err != nil {
		return "", err
	}
	if existingReportId == "" {
		fmt.Println("No existing report found, creating new one")
		newReportId, err := CreateEmptyReport(apiKey, databaseId, reportName)
		return newReportId, err
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
	switch len(existingReport.Results) {
	case 0:
		fmt.Println("No result, should create one, report ID: " + reportId)
		return "", nil
	case 1:
		fmt.Println("Result found, report ID: " + reportId)
		return existingReport.Results[0].ID, nil
	default:
		fmt.Println("Warning: too many report with same id, returning first one, report ID: " + reportId)
		return existingReport.Results[0].ID, nil
	}
}

func CreateEmptyReport(apiKey string, databaseId string, reportId string) (string, error) {
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))
	params := notion.CreatePageParams{
		ParentType: notion.ParentTypeDatabase,
		ParentID:   databaseId,

		DatabasePageProperties: &notion.DatabasePageProperties{
			"title": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: reportId,
						},
					},
				},
			},
			DATADRIFT_PROPERTY: notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: reportId,
						},
					},
				},
			},
		},
	}
	newReport, err := client.CreatePage(ctx, params)
	return newReport.ID, err
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

	propertiesToDelete := []string{}

	for _, property := range database.Properties {
		fmt.Println("Property:", property.Name)
		if property.Name == DATADRIFT_PROPERTY {
			hasDatadriftProperty = true
		}

		for _, propertyToDelete := range DefaultPropertiesToDelete {
			propertyExists := property.Name == propertyToDelete
			if propertyExists {
				propertiesToDelete = append(propertiesToDelete, propertyToDelete)
			}
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

		for _, propertyToDelete := range propertiesToDelete {
			params.Properties[propertyToDelete] = nil
		}

		fmt.Println("Creating property", params)
		updatedDB, err := client.UpdateDatabase(ctx, databaseID, params)
		if err != nil {
			return err
		}
		fmt.Println("Updated database", updatedDB, " with property", DATADRIFT_PROPERTY)
		fmt.Println("Clean empty item in database")
		queryParams := &notion.DatabaseQuery{
			Filter: &notion.DatabaseQueryFilter{
				Property: DATADRIFT_PROPERTY,
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					RichText: &notion.TextPropertyFilter{
						Equals: " ",
					},
				},
			},
		}
		emptyDatabaseItems, err := client.QueryDatabase(ctx, databaseID, queryParams)
		if err != nil {
			return err
		}
		archive := true
		for _, item := range emptyDatabaseItems.Results {
			fmt.Println("Archiving item", item.ID)
			client.UpdatePage(ctx, item.ID, notion.UpdatePageParams{
				Archived: &archive,
			})
		}
	}

	return err
}

func UpdateReport(apiKey string, reportNotionPageId string, children []notion.Block) error {
	fmt.Println("Updating report", reportNotionPageId)
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	existingReport, err := client.FindBlockChildrenByID(ctx, reportNotionPageId, &notion.PaginationQuery{PageSize: 100})
	if err != nil {
		return err
	}
	fmt.Println("Deleting children blocks:", len(existingReport.Results))

	blocks := existingReport.Results

	// Iterate over each block in existingReport.Results
	for _, block := range blocks {
		fmt.Println("Deleting block", block.ID())
		_, err := client.DeleteBlock(ctx, block.ID())
		if err != nil {
			fmt.Println("[DATADRIFT_ERROR]: deleting block", block.ID(), err)
		}
		time.Sleep(100 * time.Millisecond)

	}

	_, err = client.AppendBlockChildren(ctx, reportNotionPageId, children)
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during append", err)
	}
	return err
}
