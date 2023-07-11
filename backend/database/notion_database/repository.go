package notion_database

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/dstotijn/go-notion"
	"github.com/shopspring/decimal"
)

const PROPERTY_DATADRIFT_ID = "datadrift-id"
const PROPERTY_DATADRIFT_TIMEGRAIN = "datadrift-timegrain"
const PROPERTY_DATADRIFT_PERIOD = "datadrift-period"
const PROPERTY_DATADRIFT_DRIFT_VALUE = "datadrift-drift-value"
const PROPERTY_DATADRIFT_DIMENSION = "datadrift-dimension"

var DefaultPropertiesToDelete = []string{"Tags", "Status", "Étiquette", "Étiquettes"}

func FindOrCreateReportPageId(apiKey string, databaseId string, reportName string, period string, timeGrain common.TimeGrain, dimensionValue common.DimensionValue) (string, bool, error) {
	existingReportId, err := QueryDatabaseWithReportId(apiKey, databaseId, reportName)
	if err != nil {
		return "", false, err
	}
	if existingReportId == "" {
		fmt.Println("No existing report found, creating new one")
		newReportId, err := CreateEmptyReport(apiKey, databaseId, reportName, period, timeGrain, dimensionValue)
		return newReportId, true, err
	}
	return existingReportId, false, nil
}

func FindOrCreateSummaryReportPage(apiKey string, databaseId string, reportName string) (string, error) {
	existingReportId, err := QueryDatabaseWithReportId(apiKey, databaseId, reportName)
	if err != nil {
		return "", err
	}
	if existingReportId == "" {
		fmt.Println("No existing report found, creating new one")
		newReportId, err := CreateEmptySummaryReport(apiKey, databaseId, reportName)
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

func CreateEmptyReport(apiKey string, databaseId string, reportId string, period string, timeGrain common.TimeGrain, dimensionValue common.DimensionValue) (string, error) {
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
							Content: string(dimensionValue),
						},
					},
				},
			},
		},
	}
	newReport, err := client.CreatePage(ctx, params)
	return newReport.ID, err
}

func CreateEmptySummaryReport(apiKey string, databaseId string, reportId string) (string, error) {
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

func UpdateMetadataReport(apiKey string, reportNotionPageId string, children []notion.Block, pageProperties *notion.DatabasePageProperties) error {
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

func InitChangeLogReport(apiKey string, reportNotionPageId string, KPIInfo common.KPIReport) error {
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

	driftAmount, _ := KPIInfo.LatestValue.Sub(KPIInfo.InitialValue).Float64()

	createPageParams := notion.CreatePageParams{
		DatabasePageProperties: &notion.DatabasePageProperties{
			PROPERTY_DATADRIFT_DRIFT_VALUE: notion.DatabasePageProperty{
				Number: &driftAmount,
			},
		},
		Children: []notion.Block{
			notion.Heading1Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Overview",
						},
					},
				},
			},
			notion.ParagraphBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: KPIInfo.KPIName,
						},
						Annotations: &notion.Annotations{
							Code: true,
						},
					},
					{
						Text: &notion.Text{
							Content: " initial value was: ",
						},
					},
					{
						Text: &notion.Text{
							Content: KPIInfo.InitialValue.String(),
						},
						Annotations: &notion.Annotations{
							Bold: true,
						},
					},
				},
			},
			notion.ParagraphBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: KPIInfo.KPIName,
						},
						Annotations: &notion.Annotations{
							Code: true,
						},
					},
					{
						Text: &notion.Text{
							Content: summaryTextCurrentValueIs,
						},
					},
					{
						Text: &notion.Text{
							Content: KPIInfo.LatestValue.String(),
						},
						Annotations: &notion.Annotations{
							Bold: true,
						},
					},
				},
			},
			notion.ParagraphBlock{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: summaryTextInitialValueWas,
						},
					},
					{
						Text: &notion.Text{
							Content: displayDiff(KPIInfo.LatestValue.Sub(KPIInfo.InitialValue)),
						},
						Annotations: &notion.Annotations{
							Bold:  true,
							Color: displayDiffColor(KPIInfo.LatestValue.Sub(KPIInfo.InitialValue)),
						},
					},
				},
			},
			notion.Heading1Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Timeline",
						},
					},
				},
			},
			notion.EmbedBlock{
				URL: KPIInfo.GraphQLURL,
			},
		},
	}

	_, updateErr := client.UpdatePage(ctx, reportNotionPageId, notion.UpdatePageParams{DatabasePageProperties: *createPageParams.DatabasePageProperties})
	if updateErr != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during update", reportNotionPageId, updateErr.Error())
	}

	_, err := client.AppendBlockChildren(ctx, reportNotionPageId, createPageParams.Children)
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during append", err.Error())
	}

	reportChangeLogCreateDatabaseParams := &notion.CreateDatabaseParams{
		ParentPageID: reportNotionPageId,
		IsInline:     true,
		Title: []notion.RichText{
			{
				Text: &notion.Text{
					Content: "ChangeLog",
				},
			},
		},
		Properties: notion.DatabaseProperties{
			"Name": notion.DatabaseProperty{
				Type:  "title",
				Title: &notion.EmptyMetadata{},
			},
			"Created At": notion.DatabaseProperty{
				Type: "date",
				Date: &notion.EmptyMetadata{},
			},
			"Commit": notion.DatabaseProperty{
				Type: "url",
				URL:  &notion.EmptyMetadata{},
			},
			"Impact": notion.DatabaseProperty{
				Type: "number",
				Number: &notion.NumberMetadata{
					Format: notion.NumberFormatNumberWithCommas,
				},
			},
		},
	}
	print("\n Creating ChangeLog database...", reportChangeLogCreateDatabaseParams)
	changeLogDatabase, err := client.CreateDatabase(ctx, *reportChangeLogCreateDatabaseParams)
	if err != nil {
		fmt.Println("[DATADRIFT_ERROR]: err during changelog db creation", err.Error())
	}
	print("\n ChangeLog Database created", changeLogDatabase.ID)
	// Add all the change log to the report
	for _, event := range KPIInfo.Events {
		eventEmoji := getEventEmoji(event.Diff)
		print("\n Adding changeLog to report", event.CommitTimestamp, eventEmoji)
		_, err := client.CreatePage(ctx, notion.CreatePageParams{
			ParentID:   changeLogDatabase.ID,
			ParentType: notion.ParentTypeDatabase,
			Icon: &notion.Icon{
				Type:  notion.IconTypeEmoji,
				Emoji: &eventEmoji,
			},
			Children: []notion.Block{
				notion.ParagraphBlock{
					RichText: []notion.RichText{
						{
							Text: &notion.Text{
								Content: displayCommitComments(event),
							},
						},
					},
				},
			},
			DatabasePageProperties: &notion.DatabasePageProperties{
				"Name": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Text: &notion.Text{
								Content: displayEventTitle(event.Diff),
							},
						},
					},
				},
				"Created At": notion.DatabasePageProperty{
					Date: &notion.Date{
						Start: notion.NewDateTime(time.Unix(event.CommitTimestamp, 0), true),
					},
				},
				"Commit": notion.DatabasePageProperty{
					URL: &event.CommitUrl,
				},
				"Impact": notion.DatabasePageProperty{
					Number: &event.Diff,
				},
			},
		})
		if err != nil {
			fmt.Println("[DATADRIFT_ERROR]: err during create page", err.Error())
		}
	}

	return err
}

func displayEventTitle(diff float64) string {
	if diff == 0 {
		return "Initial Value"
	}
	return "New Drift " + displayDiff(decimal.NewFromFloat(diff))
}

func displayDiff(diff decimal.Decimal) string {
	if diff.IsPositive() {
		return "+" + diff.String()
	}
	return diff.String()
}

func displayDiffColor(diff decimal.Decimal) notion.Color {
	if diff.Equal(decimal.Zero) {
		return notion.ColorGreen
	} else if diff.IsNegative() {
		return notion.ColorOrange
	}
	return notion.ColorBlue
}

func displayCommitComments(event common.EventObject) string {
	if len(event.CommitComments) == 0 {
		return "No explanation available"
	}

	result := ""

	for _, comment := range event.CommitComments {
		result += "Author: " + comment.CommentAuthor + "\n"
		result += "Comment: " + comment.CommentBody + "\n"
		result += "\n"
	}

	if len(result) > 2000 {
		result = result[:2000]
	}
	return result
}

func getEventEmoji(diff float64) string {
	if diff == 0 {
		return "🆕"
	} else if diff > 0 {
		return "🔷"
	}
	return "🔶"
}

const summaryTextInitialValueWas = "Total drift since initial value: "
const summaryTextCurrentValueIs = " current value is: "
