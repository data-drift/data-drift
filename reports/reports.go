package reports

import (
	"fmt"
	"strconv"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	notion_database "github.com/data-drift/kpi-git-history/database/notion"
	"github.com/dstotijn/go-notion"
)

func CreateReport(syncConfig common.SyncConfig, KPIInfo common.KPIReport) error {
	reportNotionPageId, _ := notion_database.FindOrCreateReportPageId(syncConfig.NotionAPIKey, syncConfig.NotionDatabaseID, KPIInfo.KPIName)
	fmt.Println(reportNotionPageId)

	params := notion.CreatePageParams{
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
							Content: strconv.FormatFloat(KPIInfo.InitialValue, 'f', 2, 64),
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
							Content: strconv.FormatFloat(KPIInfo.LatestValue, 'f', 2, 64),
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
							Content: displayDiff(KPIInfo.LatestValue - KPIInfo.InitialValue),
						},
						Annotations: &notion.Annotations{
							Bold:  true,
							Color: displayDiffColor(KPIInfo.LatestValue - KPIInfo.InitialValue),
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
						Content: strconv.FormatFloat(KPIInfo.InitialValue, 'f', 2, 64),
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
						Content: displayDiff(event.Diff),
					},
					Annotations: &notion.Annotations{
						Bold:  true,
						Color: displayDiffColor(event.Diff),
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
		if event.EventType == "create" {
			children = append(children, driftEventDate, bulletListFirstItemCreateEvent)
		} else {
			children = append(children, driftEventDate, bulletListFirstItemUpdateEvent, bulletListSecondItemUpdateEvent)
		}
	}
	params.Children = append(params.Children, children...)

	err := notion_database.UpdateReport(syncConfig.NotionAPIKey, reportNotionPageId, params.Children)
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err)
	}

	return nil
}

func displayDiff(diff float64) string {
	if diff >= 0 {

		return "+" + strconv.FormatFloat(diff, 'f', 2, 64)
	}
	return strconv.FormatFloat(diff, 'f', 2, 64)
}

func displayDiffColor(diff float64) notion.Color {
	if diff < 0 {
		return notion.ColorOrange
	}
	return notion.ColorBlue
}

func GetTimeGrain(periodKey string) (common.TimeGrain, error) {
	_, err := time.Parse("2006-01-02", periodKey)
	if err == nil {
		return common.Day, nil
	}
	_, err = time.Parse("2006-W01", fmt.Sprintf("%s-01", periodKey))
	if err == nil {
		return common.Week, nil
	}
	_, err = time.Parse("2006-01", fmt.Sprintf("%s-01", periodKey))
	if err == nil {
		return common.Month, nil
	}
	_, err = time.Parse("2006", fmt.Sprintf("%s-01-01", periodKey))
	if err == nil {
		return common.Year, nil
	}
	return "", fmt.Errorf("invalid period key: %s", periodKey)
}
