package reports

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/data-drift/data-drift/common"
	"github.com/data-drift/data-drift/database/notion_database"
	"github.com/dstotijn/go-notion"
	"github.com/shopspring/decimal"
)

func CreateReport(syncConfig common.SyncConfig, KPIInfo common.KPIReport) error {
	timeGrain, _ := GetTimeGrain(KPIInfo.PeriodId)
	reportNotionPageId, findOrCreateError := notion_database.FindOrCreateReportPageId(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, KPIInfo.KPIName, string(KPIInfo.PeriodId), timeGrain, KPIInfo.DimensionValue)
	if findOrCreateError != nil {
		return fmt.Errorf("failed to create reportNotionPageId: %v", findOrCreateError.Error())
	}

	driftAmount, _ := KPIInfo.LatestValue.Sub(KPIInfo.InitialValue).Float64()

	createPageParams := notion.CreatePageParams{
		DatabasePageProperties: &notion.DatabasePageProperties{
			notion_database.PROPERTY_DATADRIFT_DRIFT_VALUE: notion.DatabasePageProperty{
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
							Content: " current value is: ",
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
							Content: "Total drift since initial value: ",
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
			notion.Heading1Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Changelog",
						},
					},
				},
			},
		},
	}
	var children []notion.Block
	for _, event := range KPIInfo.Events {
		driftEventDate := notion.ParagraphBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "🗓 Event ",
					},
				},

				{
					Mention: &notion.Mention{
						Type: notion.MentionTypeDate,
						Date: &notion.Date{
							Start: notion.NewDateTime(time.Unix(event.CommitTimestamp, 0), true),
						},
					},
				},
			},
		}
		bulletListFirstItemCreateEvent := notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Initial value: ",
					},
				},
				{
					Text: &notion.Text{
						Content: KPIInfo.InitialValue.String(),
					},
					Annotations: &notion.Annotations{
						Bold:  true,
						Color: notion.ColorGray,
					},
				},
			},
		}
		bulletListFirstItemUpdateEvent := notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Impact: ",
					},
				},
				{
					Text: &notion.Text{
						Content: displayDiff(decimal.NewFromFloat(event.Diff)),
					},
					Annotations: &notion.Annotations{
						Bold:  true,
						Color: displayDiffColor(decimal.NewFromFloat(event.Diff)),
					},
				},
			},
		}
		bulletListSecondItemUpdateEvent := notion.BulletedListItemBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "commit",
						Link:    &notion.Link{URL: event.CommitUrl},
					},
				},
			},
		}
		toggleUpdateEvent := notion.ToggleBlock{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Explanation",
					},
				},
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
		}
		if event.EventType == "create" {
			children = append(children, driftEventDate, bulletListFirstItemCreateEvent)
		} else {
			children = append(children, driftEventDate, bulletListFirstItemUpdateEvent, bulletListSecondItemUpdateEvent, toggleUpdateEvent)
		}
	}
	createPageParams.Children = append(createPageParams.Children, children...)

	err := notion_database.UpdateReport(syncConfig.NotionAPIKey, reportNotionPageId, createPageParams.Children, createPageParams.DatabasePageProperties)
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err.Error())
	}

	return nil
}

func CreateSummaryReport(syncConfig common.SyncConfig, metricConfig common.MetricConfig, chartUrls map[common.TimeGrain]string) error {
	fmt.Println("Creating summary report")
	reportNotionPageId, findOrCreateError := notion_database.FindOrCreateSummaryReportPage(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, "Summary of "+metricConfig.MetricName)
	fmt.Println(reportNotionPageId)
	fmt.Println(findOrCreateError)

	var children []notion.Block
	for _, timeGrain := range []common.TimeGrain{common.Day, common.Week, common.Month, common.Quarter, common.Year} {
		chartUrl := chartUrls[timeGrain]
		if chartUrl == "" {
			continue
		}
		children = append(children, notion.Heading1Block{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: "Cohorts " + string(timeGrain),
					},
				},
			},
		},
			notion.EmbedBlock{
				URL: chartUrl,
			},
		)

	}
	err := notion_database.UpdateReport(syncConfig.NotionAPIKey, reportNotionPageId, children, &notion.DatabasePageProperties{})
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err.Error())
	}

	return nil
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

func GetTimeGrain(periodKeyParam common.PeriodKey) (common.TimeGrain, error) {
	periodKey := string(periodKeyParam)
	_, err := time.Parse("2006-01-02", periodKey)
	if err == nil {
		return common.Day, nil
	}
	_, err = ParseYearWeek(periodKey)
	if err == nil {
		return common.Week, nil
	}
	_, err = time.Parse("2006-01", periodKey)
	if err == nil {
		return common.Month, nil
	}
	_, err = ParseQuarterDate(periodKey)
	if err == nil {
		return common.Quarter, nil
	}
	_, err = time.Parse("2006", periodKey)
	if err == nil {
		return common.Year, nil
	}
	return "", fmt.Errorf("invalid period key: %s", periodKey)
}

func ParseYearWeek(yearWeek string) (time.Time, error) {
	if len(yearWeek) != 8 {
		return time.Time{}, fmt.Errorf("invalid year week format: %s", yearWeek)
	}
	year, err := strconv.Atoi(yearWeek[0:4])
	if err != nil {
		return time.Time{}, err
	}

	week, err := strconv.Atoi(yearWeek[6:])
	if err != nil {
		return time.Time{}, err
	}

	// Get the first day of the week (Monday)
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 7*(week-1)+1)

	return firstDay, nil
}

func ParseQuarterDate(s string) (time.Time, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid quarter date format: %s", s)
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year format in quarter date: %s", s)
	}
	quarter := parts[1]
	switch quarter {
	case "Q1":
		return time.Date(year, time.March, 31, 0, 0, 0, 0, time.UTC), nil
	case "Q2":
		return time.Date(year, time.June, 30, 0, 0, 0, 0, time.UTC), nil
	case "Q3":
		return time.Date(year, time.September, 30, 0, 0, 0, 0, time.UTC), nil
	case "Q4":
		return time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC), nil
	default:
		return time.Time{}, fmt.Errorf("invalid quarter format in quarter date: %s", s)
	}
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
