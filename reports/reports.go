package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/dstotijn/go-notion"
	"github.com/sanity-io/litter"
)

type httpTransport struct {
	w io.Writer
}

// RoundTrip implements http.RoundTripper. It multiplexes the read HTTP response
// data to an io.Writer for debugging.
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}

func CreateReport(syncConfig common.SyncConfig, KPIInfo common.KPIInfo) error {
	fmt.Println("CreateReport called with", KPIInfo)
	ctx := context.Background()
	apiKey := syncConfig.NotionAPIKey
	databaseId := syncConfig.NotionDatabaseID

	if apiKey == "" || databaseId == "" {
		return fmt.Errorf("missing Notion API key or database ID")
	}
	buf := &bytes.Buffer{}
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
							Content: KPIInfo.KPIName,
						},
					},
				},
			},
		},

		Children: []notion.Block{
			notion.EmbedBlock{
				URL: KPIInfo.GraphQLURL,
			},
		},
	}
	page, err := client.CreatePage(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create page: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		return fmt.Errorf("failed to decode result: %v", err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		return fmt.Errorf("failed to decode result: %v", err)
	}

	// Pretty print parsed `notion.Page` value.
	litter.Dump(page.ID)
	return nil
}
