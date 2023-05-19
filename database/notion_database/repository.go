package notion_database

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/dstotijn/go-notion"
)

const PROPERTY_DATADRIFT_ID = "datadrift-id"
const PROPERTY_DATADRIFT_TIMEGRAIN = "datadrift-timegrain"
const PROPERTY_DATADRIFT_PERIOD = "datadrift-period"
const PROPERTY_DATADRIFT_DRIFT_VALUE = "datadrift-drift-value"
const PROPERTY_DATADRIFT_DIMENSION = "datadrift-dimension"

var DefaultPropertiesToDelete = []string{"Tags", "Status", "Étiquette", "Étiquettes"}

func FindOrCreateReportPageId(apiKey string, databaseId string, reportName string, period string, timeGrain common.TimeGrain) (string, error) {
	existingReportId, err := QueryDatabaseWithReportId(apiKey, databaseId, reportName)
	if err != nil {
		return "", err
	}
	if existingReportId == "" {
		fmt.Println("No existing report found, creating new one")
		newReportId, err := CreateEmptyReport(apiKey, databaseId, reportName, period, timeGrain)
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

func CreateEmptyReport(apiKey string, databaseId string, reportId string, period string, timeGrain common.TimeGrain) (string, error) {
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
			PROPERTY_DATADRIFT_ID: notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: reportId,
						},
					},
				},
			},
			PROPERTY_DATADRIFT_PERIOD: notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: period,
						},
					},
				},
			},
			PROPERTY_DATADRIFT_TIMEGRAIN: notion.DatabasePageProperty{
				Select: &notion.SelectOptions{
					Name: string(timeGrain),
				},
			},
			PROPERTY_DATADRIFT_DIMENSION: notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "coucou",
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

	shouldCreateDatadriftPropertyId := true
	shouldCreateDatadriftPropertyPeriod := true
	shouldCreateDatadriftPropertyTimeGrain := true
	shouldCreateDatadriftPropertyDriftValue := true
	shouldCreateDatadriftPropertyDimension := true

	propertiesToDelete := []string{}

	for _, property := range database.Properties {
		fmt.Println("Property:", property.Name)
		if property.Name == PROPERTY_DATADRIFT_ID {
			shouldCreateDatadriftPropertyId = false
		}
		if property.Name == PROPERTY_DATADRIFT_PERIOD {
			shouldCreateDatadriftPropertyPeriod = false
		}
		if property.Name == PROPERTY_DATADRIFT_TIMEGRAIN {
			shouldCreateDatadriftPropertyTimeGrain = false
		}
		if property.Name == PROPERTY_DATADRIFT_DRIFT_VALUE {
			shouldCreateDatadriftPropertyDriftValue = false
		}
		if property.Name == PROPERTY_DATADRIFT_DIMENSION {
			shouldCreateDatadriftPropertyDimension = false
		}

		for _, propertyToDelete := range DefaultPropertiesToDelete {
			propertyExists := property.Name == propertyToDelete
			if propertyExists {
				propertiesToDelete = append(propertiesToDelete, propertyToDelete)
			}
		}

	}
	fmt.Println("hasDatadriftProperty:", shouldCreateDatadriftPropertyId)
	shouldCreateProperties := shouldCreateDatadriftPropertyId || shouldCreateDatadriftPropertyPeriod || shouldCreateDatadriftPropertyTimeGrain || shouldCreateDatadriftPropertyDriftValue || shouldCreateDatadriftPropertyDimension
	if shouldCreateProperties {
		params := notion.UpdateDatabaseParams{
			Properties: map[string]*notion.DatabaseProperty{},
		}

		for _, propertyToDelete := range propertiesToDelete {
			params.Properties[propertyToDelete] = nil
		}

		if shouldCreateDatadriftPropertyId {
			params.Properties[PROPERTY_DATADRIFT_ID] = &notion.DatabaseProperty{
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			}
		}

		if shouldCreateDatadriftPropertyPeriod {
			params.Properties[PROPERTY_DATADRIFT_PERIOD] = &notion.DatabaseProperty{
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			}
		}

		if shouldCreateDatadriftPropertyDimension {
			params.Properties[PROPERTY_DATADRIFT_DIMENSION] = &notion.DatabaseProperty{
				Type:     notion.DBPropTypeRichText,
				RichText: &notion.EmptyMetadata{},
			}
		}

		if shouldCreateDatadriftPropertyTimeGrain {
			params.Properties[PROPERTY_DATADRIFT_TIMEGRAIN] = &notion.DatabaseProperty{
				Type: notion.DBPropTypeSelect,
				Select: &notion.SelectMetadata{
					Options: []notion.SelectOptions{
						{Name: string(common.Day), Color: notion.ColorYellow},
						{Name: string(common.Month), Color: notion.ColorOrange},
						{Name: string(common.Week), Color: notion.ColorRed},
						{Name: string(common.Quarter), Color: notion.ColorPink},
						{Name: string(common.Year), Color: notion.ColorPurple},
					},
				},
			}
		}

		if shouldCreateDatadriftPropertyDriftValue {
			params.Properties[PROPERTY_DATADRIFT_DRIFT_VALUE] = &notion.DatabaseProperty{
				Type: notion.DBPropTypeNumber,
				Number: &notion.NumberMetadata{
					Format: notion.NumberFormatNumberWithCommas,
				},
			}
		}

		fmt.Println("Creating property", params)
		_, err := client.UpdateDatabase(ctx, databaseID, params)
		if err != nil {
			return err
		}
		fmt.Println("Clean empty item in database")
		queryParams := &notion.DatabaseQuery{
			Filter: &notion.DatabaseQueryFilter{
				Property: PROPERTY_DATADRIFT_ID,
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

func UpdateReport(apiKey string, reportNotionPageId string, children []notion.Block, pageProperties *notion.DatabasePageProperties) error {
	if reportNotionPageId == "" {
		fmt.Println("No report page id provided")
		return nil
	}

	fmt.Println("Updating report", reportNotionPageId)
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	_, updateErr := client.UpdatePage(ctx, reportNotionPageId, notion.UpdatePageParams{DatabasePageProperties: *pageProperties})
	if updateErr != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during update", reportNotionPageId, updateErr.Error())
	}
	existingReport, err := client.FindBlockChildrenByID(ctx, reportNotionPageId, &notion.PaginationQuery{PageSize: 100})
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during find block", reportNotionPageId, err.Error())
		return err
	}
	fmt.Println("Deleting children blocks:", len(existingReport.Results))

	blocks := existingReport.Results

	// Iterate over each block in existingReport.Results
	for _, block := range blocks {
		fmt.Println("Deleting block", block.ID())
		_, err := client.DeleteBlock(ctx, block.ID())
		if err != nil {
			fmt.Println("[DATADRIFT_ERROR]: deleting block", block.ID(), err.Error())
		}
		time.Sleep(100 * time.Millisecond)

	}

	_, err = client.AppendBlockChildren(ctx, reportNotionPageId, children)
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during append", err.Error())
	}
	return err
}
