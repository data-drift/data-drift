package reports

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/data-drift/kpi-git-history/common"
	"github.com/dstotijn/go-notion"
	"github.com/joho/godotenv"
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

func CreateReport(KPIInfo common.KPIInfo) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	apiKey := os.Getenv("NOTION_API_KEY")
	databaseId := os.Getenv("NOTION_DATABASE_ID")
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
		log.Fatalf("Failed to create page: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print parsed `notion.Page` value.
	litter.Dump(page.ID)
}
