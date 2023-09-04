package notion_database

import (
	"testing"

	"github.com/data-drift/data-drift/common"
	"github.com/go-playground/assert/v2"
	"github.com/shopspring/decimal"
)

func TestDisplayEventTitle(t *testing.T) {
	// Test the "create" event type
	createEvent := common.EventObject{
		EventType: common.EventTypeCreate,
		Current:   decimal.NewFromFloat(10.1),
	}
	assert.Equal(t, "Initial Value 10.1", displayEventTitle(createEvent))

	// Test the "update" event type
	updateEvent := common.EventObject{
		EventType: common.EventTypeUpdate,
		Diff:      2.5,
	}
	assert.Equal(t, "New Drift +2.5", displayEventTitle(updateEvent))

	// Test the "delete" event type
	deleteEvent := common.EventObject{
		EventType: common.EventTypeUpdate,
		Diff:      -1.75,
	}
	assert.Equal(t, "New Drift -1.75", displayEventTitle(deleteEvent))
}
