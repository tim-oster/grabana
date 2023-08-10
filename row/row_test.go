package row

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestNewRowsCanBeCreated(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "Some row")

	req.NoError(err)
	req.Equal("Some row", panel.builder.Title)
	req.True(panel.builder.ShowTitle)
}

func TestRowsCanHaveHiddenTitle(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", HideTitle())

	req.NoError(err)
	req.False(panel.builder.ShowTitle)
}

func TestRowsCanHaveVisibleTitle(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", ShowTitle())

	req.NoError(err)
	req.True(panel.builder.ShowTitle)
}

func TestRowsCanHaveGraphs(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithGraph("HTTP Rate"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveTimeSeries(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithTimeSeries("HTTP Rate"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

//func TestRowsCanHaveTimeSeriesAndAlert(t *testing.T) {
//	req := require.New(t)
//	board := sdk.NewBoard("")
//
//	panel, err := New(
//		board,
//		"",
//		WithTimeSeries(
//			"HTTP Rate",
//			timeseries.DataSource("Prometheus"),
//			timeseries.Alert("Too many heap allocations"),
//		),
//	)
//
//	req.NoError(err)
//	req.Len(panel.builder.Panels, 1)
//	req.Len(panel.Alerts, 1)
//
//	req.Equal("Too many heap allocations", panel.Alerts[0].Builder.Title)
//}

func TestRowsCanHaveTextPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithText("HTTP Rate"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveTablePanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithTable("Some table"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveSingleStatPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithSingleStat("Some stat"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveStatPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithStat("Some stat"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveHeatmapPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithHeatmap("Some heatmap"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveLogsPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithLogs("Some logs"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveGaugePanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", WithGauge("Some gauge"))

	req.NoError(err)
	req.Len(panel.builder.Panels, 1)
}

func TestRowsCanHaveRepeatedPanels(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", RepeatFor("repeated"))

	req.NoError(err)
	req.Equal("repeated", *panel.builder.Repeat)
}

func TestRowsCanBeCollapsedByDefault(t *testing.T) {
	req := require.New(t)
	board := sdk.NewBoard("")

	panel, err := New(board, "", Collapse())

	req.NoError(err)
	req.True(panel.builder.Collapse)
}
