package extractor

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/rodrwan/news-getter/graph/model"
	"gopkg.in/yaml.v2"
)

type Extractor struct {
	Path    string
	Sources map[string]*Source
}

type Source struct {
	URL      string `yaml:"url"`
	Headline string `yaml:"headline"`
	Link     string `yaml:"link"`

	Body io.ReadCloser
}

func (e *Extractor) Load(countries ...string) error {
	sources := make(map[string]*Source, len(countries))

	for _, country := range countries {
		var s Source
		file := fmt.Sprintf("%s/%s.yaml", e.Path, country)
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.Wrap(err, "yamlFile.Get error")
		}

		if err := yaml.Unmarshal(yamlFile, &s); err != nil {
			return errors.Wrap(err, "Unmarshal")
		}

		sources[country] = &s
	}

	e.Sources = sources
	return nil
}

func (e *Extractor) GetHTML(ctx context.Context) error {
	for _, source := range e.Sources {
		client := &http.Client{}

		req, err := http.NewRequestWithContext(ctx, "GET", source.URL, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		source.Body = resp.Body
	}

	return nil
}

func (e *Extractor) GetNews() ([]*model.NewsItem, error) {
	// Load the HTML document
	reg, err := regexp.Compile("[\\n\\s]+")
	if err != nil {
		return nil, err
	}

	var news []*model.NewsItem
	for _, source := range e.Sources {
		doc, err := goquery.NewDocumentFromReader(source.Body)
		if err != nil {
			return nil, err
		}
		n := &model.NewsItem{}

		// Find the news
		doc.Find(source.Link).Each(func(i int, s *goquery.Selection) {
			if link, ok := s.Attr("href"); ok {
				n.Link = link
			}
		})

		doc.Find(source.Headline).Each(func(i int, s *goquery.Selection) {
			if s != nil {
				n.Headline = strings.TrimSpace(reg.ReplaceAllString(s.Text(), " "))
			}
		})
		news = append(news, n)

		defer source.Body.Close()
	}

	return news, nil
}
