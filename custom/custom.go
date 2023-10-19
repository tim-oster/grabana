package custom

import (
	"fmt"
	alert "github.com/K-Phoen/grabana/ngalert"
	"github.com/K-Phoen/sdk"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(cst *Custom) error

// Custom represents a custom panel.
type Custom struct {
	Builder *sdk.Panel
	Alerts  []*alert.Alert
}

// New creates a new custom panel.
func New(title string, options ...Option) (*Custom, error) {
	panel := &Custom{Builder: sdk.NewCustom(title)}
	panel.Builder.IsNew = false

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
		Span(6),
	}
}

// CustomConfig sets the custom panel config of this panel.
func CustomConfig(config map[string]any) Option {
	return func(cst *Custom) error {
		*cst.Builder.CustomPanel = config
		return nil
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(cst *Custom) error {
		cst.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			cst.Builder.Links = append(cst.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the graph.
func DataSource(source string) Option {
	return func(cst *Custom) error {
		cst.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(cst *Custom) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		cst.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(cst *Custom) error {
		cst.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(cst *Custom) error {
		cst.Builder.Description = &content

		return nil
	}
}

// Alert creates a next generation alert (grafana unified alerting) for this graph.
func Alert(name string, opts ...alert.Option) Option {
	return func(cst *Custom) error {
		obj := alert.New(name, opts...)

		for i, data := range obj.Builder.Data {
			if data.DatasourceUid == "" {
				data.DatasourceUid = cst.Builder.Datasource.LegacyName
				obj.Builder.Data[i] = data
			}
		}

		cst.Alerts = append(cst.Alerts, obj)

		return nil
	}
}
