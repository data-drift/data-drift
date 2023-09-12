package notion_database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/helpers"
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
		log.Println("No existing report found, creating new one")
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
		log.Println("No existing report found, creating new one")
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
		log.Println("No result, should create one, report ID: " + reportId)
		return "", nil
	case 1:
		log.Println("Result found, report ID: " + reportId)
		return existingReport.Results[0].ID, nil
	default:
		log.Println("Warning: too many report with same id, returning first one, report ID: " + reportId)
		return existingReport.Results[0].ID, nil
	}
}

func QueryDatabaseWithMetricAndTimegrain(apiKey string, databaseId string, metricName string, timeGrain common.TimeGrain) ([]notion.Page, error) {
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	queryParams := &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: []notion.DatabaseQueryFilter{
				{
					Property: "datadrift-id",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						RichText: &notion.TextPropertyFilter{
							StartsWith: metricName,
						},
					},
				},
				{
					Property: "datadrift-timegrain",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Select: &notion.SelectDatabaseQueryFilter{
							Equals: string(timeGrain),
						},
					},
				},
			},
		},
		Sorts: []notion.DatabaseQuerySort{
			{Property: "datadrift-id", Direction: "ascending"},
		},
	}

	existingReport, err := client.QueryDatabase(ctx, databaseId, queryParams)
	return existingReport.Results, err
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
		log.Println("Property:", property.Name)
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
	log.Println("hasDatadriftProperty:", shouldCreateDatadriftPropertyId)
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

		log.Println("Creating property", params)
		_, err := client.UpdateDatabase(ctx, databaseID, params)
		if err != nil {
			return err
		}
		log.Println("Clean empty item in database")
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
			log.Println("Archiving item", item.ID)
			client.UpdatePage(ctx, item.ID, notion.UpdatePageParams{
				Archived: &archive,
			})
		}
	}

	return err
}

func UpdateMetadataReport(apiKey string, reportNotionPageId string, children []notion.Block, pageProperties *notion.DatabasePageProperties) error {
	if reportNotionPageId == "" {
		log.Println("No report page id provided")
		return nil
	}

	log.Println("Updating report", reportNotionPageId)
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	_, updateErr := client.UpdatePage(ctx, reportNotionPageId, notion.UpdatePageParams{DatabasePageProperties: *pageProperties})
	if updateErr != nil {
		log.Println("[DATADRIFT_ERROR]: err during update", reportNotionPageId, updateErr.Error())
	}
	existingReport, err := client.FindBlockChildrenByID(ctx, reportNotionPageId, &notion.PaginationQuery{PageSize: 100})
	if err != nil {
		log.Println("[DATADRIFT_ERROR]: err during find block", reportNotionPageId, err.Error())
		return err
	}
	log.Println("Deleting children blocks:", len(existingReport.Results))

	blocks := existingReport.Results

	// Iterate over each block in existingReport.Results
	for _, block := range blocks {
		log.Println("Deleting block", block.ID())
		_, err := client.DeleteBlock(ctx, block.ID())
		if err != nil {
			log.Println("[DATADRIFT_ERROR]: deleting block", block.ID(), err.Error())
		}
		time.Sleep(100 * time.Millisecond)

	}

	_, err = client.AppendBlockChildren(ctx, reportNotionPageId, children)
	if err != nil {
		log.Println("[DATADRIFT_ERROR]: err during append", err.Error())
	}
	return err
}

func InitChangeLogReport(apiKey string, reportNotionPageId string, KPIInfo common.KPIReport) error {
	if reportNotionPageId == "" {
		log.Println("No report page id provided")
		return nil
	}

	log.Println("Updating report", reportNotionPageId)
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	updatePageProperties(ctx, KPIInfo, client, reportNotionPageId)

	createPageChildren := []notion.Block{
		notion.Heading1Block{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Overview",
					},
				},
			},
		},
		buildInitialValueParagraph(KPIInfo),
		buildCurrentValueParagraph(KPIInfo),
		buildDriftParagraph(KPIInfo),
		notion.Heading1Block{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Timeline",
					},
				},
			},
		},
		buildEmberChartBlock(KPIInfo),
	}

	_, err := client.AppendBlockChildren(ctx, reportNotionPageId, createPageChildren)
	if err != nil {
		log.Println("[DATADRIFT_ERROR]: err during append", err.Error())
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
			"datadrift-event-id": notion.DatabaseProperty{
				Type:     "rich_text",
				RichText: &notion.EmptyMetadata{},
			},
			"Impact": notion.DatabaseProperty{
				Type: "number",
				Number: &notion.NumberMetadata{
					Format: notion.NumberFormatNumberWithCommas,
				},
			},
		},
	}
	log.Println("Creating ChangeLog database...", reportChangeLogCreateDatabaseParams)
	changeLogDatabase, err := client.CreateDatabase(ctx, *reportChangeLogCreateDatabaseParams)
	if err != nil {
		log.Println("[DATADRIFT_ERROR]: err during changelog db creation", err.Error())
	}
	log.Println("ChangeLog Database created", changeLogDatabase.ID)
	// Add all the change log to the report
	for _, event := range KPIInfo.Events {
		err := createEventInNotionReport(event, client, ctx, changeLogDatabase.ID)
		if err != nil {
			log.Println("[DATADRIFT_ERROR]: err during create page", err.Error())
		}
	}

	return err
}

func updatePageProperties(ctx context.Context, KPIInfo common.KPIReport, client *notion.Client, reportNotionPageId string) {
	driftAmount, _ := KPIInfo.LatestValue.Sub(KPIInfo.InitialValue).Float64()

	updatePageProperties := notion.DatabasePageProperties{
		PROPERTY_DATADRIFT_DRIFT_VALUE: notion.DatabasePageProperty{
			Number: &driftAmount,
		},
	}

	_, updateErr := client.UpdatePage(ctx, reportNotionPageId, notion.UpdatePageParams{DatabasePageProperties: updatePageProperties})
	if updateErr != nil {
		log.Println("[DATADRIFT_ERROR]: err during update", reportNotionPageId, updateErr.Error())
	}
}

func createEventInNotionReport(event common.EventObject, client *notion.Client, ctx context.Context, changeLogDatabaseId string) error {
	eventEmoji := getEventEmoji(event.Diff)
	log.Println("Adding changeLog to report", event.CommitTimestamp, eventEmoji)
	_, err := client.CreatePage(ctx, notion.CreatePageParams{
		ParentID:   changeLogDatabaseId,
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
							Content: displayEventTitle(event),
						},
					},
				},
			},
			"Created At": notion.DatabasePageProperty{
				Date: &notion.Date{
					Start: notion.NewDateTime(time.Unix(event.CommitTimestamp, 0), true),
				},
			},
			"Impact": notion.DatabasePageProperty{
				Number: &event.Diff,
			},
			"Commit": notion.DatabasePageProperty{
				URL: &event.CommitUrl,
			},
			"datadrift-event-id": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: generateEventId(event),
						},
						Type: "text",
					},
				},
			},
		},
	})
	return err
}

func UpdateChangeLogReport(apiKey string, reportNotionPageId string, KPIInfo common.KPIReport) error {
	if reportNotionPageId == "" {
		log.Println("No report page id provided")
		return nil
	}

	log.Println("Updating report", reportNotionPageId)
	buf := &bytes.Buffer{}
	ctx := context.Background()

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	updatePageProperties(ctx, KPIInfo, client, reportNotionPageId)

	pageContent, err := client.FindBlockChildrenByID(ctx, reportNotionPageId, &notion.PaginationQuery{PageSize: 100})
	if err != nil {
		return err
	}

	var driftBlock *notion.ParagraphBlock
	var currentValueBlock *notion.ParagraphBlock
	var embedChartBlock *notion.EmbedBlock
	var changeLogDatabaseId string
	for _, block := range pageContent.Results {
		switch b := block.(type) {
		case *notion.ParagraphBlock:
			if len(b.RichText) > 0 && b.RichText[0].Text.Content == summaryTextInitialValueWas {
				driftBlock = b
				break
			}
			if len(b.RichText) > 1 && b.RichText[1].Text.Content == summaryTextCurrentValueIs {
				currentValueBlock = b
				break
			}
		case *notion.EmbedBlock:
			log.Println("block is an embed block", strings.HasPrefix(b.URL, "https://app.data-drift.io/report"), b.URL)
			if strings.HasPrefix(b.URL, "https://app.data-drift.io/report") {
				embedChartBlock = b
			}
		case *notion.ChildDatabaseBlock:
			log.Println("block is a child database block", b.ID())
			changeLogDatabaseId = b.ID()
		default:
			log.Println("block is not a known block type")
		}
	}
	if changeLogDatabaseId != "" {
		log.Println("Adding missing events in ChangeLog database...", changeLogDatabaseId)
		createMissingEvents(client, ctx, changeLogDatabaseId, KPIInfo)
	}
	if driftBlock != nil {
		log.Println("Updating driftBlock: ", driftBlock.ID())
		blockID := driftBlock.ID()
		newContent := buildDriftParagraph(KPIInfo)

		_, err := client.UpdateBlock(ctx, blockID, newContent)
		if err != nil {
			log.Println("Error updating driftBlock: ", err.Error())
		}
	}

	if currentValueBlock != nil {
		log.Println("Updating currentValueBlock: ", currentValueBlock.ID())
		blockID := currentValueBlock.ID()
		newContent := buildCurrentValueParagraph(KPIInfo)

		_, err := client.UpdateBlock(ctx, blockID, newContent)
		if err != nil {
			log.Println("Error updating currentValueBlock: ", err.Error())
		}
	}
	if embedChartBlock != nil {
		log.Println("Updating embedChartBlock: ", embedChartBlock.ID())
		blockID := embedChartBlock.ID()
		newContent := buildEmberChartBlock(KPIInfo)

		_, err := client.UpdateBlock(ctx, blockID, newContent)
		if err != nil {
			log.Println("Error updating embedChartBlock: ", err.Error())
		}
	}

	return nil
}

func createMissingEvents(client *notion.Client, ctx context.Context, databaseID string, KPIInfo common.KPIReport) error {
	// Get the database
	db, err := client.QueryDatabase(ctx, databaseID, nil)
	if err != nil {
		return err
	}

	eventsToCreate := make(map[string]bool)
	for _, event := range KPIInfo.Events {
		log.Println("Adding event to create  ", event.CommitTimestamp)
		eventsToCreate[generateEventId(event)] = true
	}

	type ChangeLogRichText struct {
		PlainText string `json:"plain_text"`
	}

	type ChangeLogDatadriftID struct {
		RichText []ChangeLogRichText `json:"rich_text"`
	}

	type ChangeLogProperties struct {
		DatadriftID ChangeLogDatadriftID `json:"datadrift-event-id"`
	}

	// Get the existing events
	for _, page := range db.Results {
		// Type-assert Properties to a map[string]notion.PropertyValue
		properties := page.Properties

		jsonProperties, _ := json.Marshal(properties)

		var propertiesMap ChangeLogProperties
		if err := json.Unmarshal(jsonProperties, &propertiesMap); err != nil {
			log.Println(err.Error())
		}

		// Access the "Commit" property

		changeLogEventId := propertiesMap.DatadriftID.RichText[0].PlainText
		eventsToCreate[changeLogEventId] = false

	}

	log.Println(eventsToCreate)
	for _, event := range KPIInfo.Events {
		if eventsToCreate[generateEventId(event)] {
			log.Println("Creating the event: ", generateEventId(event))
			err := createEventInNotionReport(event, client, ctx, databaseID)
			if err != nil {
				log.Println("[DATADRIFT_ERROR]: err during create page", err.Error())
			}
		} else {
			log.Println("Event already exist: ", generateEventId(event))
		}
	}

	return nil
}

func displayEventTitle(event common.EventObject) string {
	if event.EventType == common.EventTypeCreate {
		return fmt.Sprintf("Initial Value %s", helpers.FormatWithSeparator(event.Current))
	}
	return "New Drift " + displayDiff(decimal.NewFromFloat(event.Diff))
}

func displayDiff(diff decimal.Decimal) string {
	if diff.IsPositive() {
		return "+" + helpers.FormatWithSeparator(diff)
	}
	return helpers.FormatWithSeparator(diff)
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

func buildCurrentValueParagraph(KPIInfo common.KPIReport) notion.ParagraphBlock {
	return notion.ParagraphBlock{
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
					Content: helpers.FormatWithSeparator(KPIInfo.LatestValue),
				},
				Annotations: &notion.Annotations{
					Bold: true,
				},
			},
		},
	}
}

func buildInitialValueParagraph(KPIInfo common.KPIReport) notion.ParagraphBlock {
	return notion.ParagraphBlock{
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
					Content: helpers.FormatWithSeparator(KPIInfo.InitialValue),
				},
				Annotations: &notion.Annotations{
					Bold: true,
				},
			},
		},
	}
}

func buildDriftParagraph(KPIInfo common.KPIReport) notion.ParagraphBlock {
	return notion.ParagraphBlock{
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
	}
}

func buildEmberChartBlock(KPIInfo common.KPIReport) notion.EmbedBlock {
	return notion.EmbedBlock{
		URL: KPIInfo.WaterfallChartUrl,
	}
}

func generateEventId(event common.EventObject) string {
	return fmt.Sprintf("event-timestamp-%d", event.CommitTimestamp)
}
