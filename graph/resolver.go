package graph

import "github.com/rodrwan/news-getter/services/extractor"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Extractor *extractor.Extractor
}
