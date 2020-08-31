package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/rodrwan/news-getter/graph/generated"
	"github.com/rodrwan/news-getter/graph/model"
)

func (r *queryResolver) GetNewsByCountry(ctx context.Context, country string) (*model.NewsItem, error) {
	if err := r.Resolver.Extractor.Load(country); err != nil {
		return nil, err
	}

	if err := r.Resolver.Extractor.GetHTML(ctx); err != nil {
		return nil, err
	}

	news, err := r.Resolver.Extractor.GetNews()
	if err != nil {
		return nil, err
	}

	return news[0], nil
}

func (r *queryResolver) GetNews(ctx context.Context) ([]*model.NewsItem, error) {
	if err := r.Resolver.Extractor.Load("cl", "es"); err != nil {
		return nil, err
	}

	if err := r.Resolver.Extractor.GetHTML(ctx); err != nil {
		return nil, err
	}

	return r.Resolver.Extractor.GetNews()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
