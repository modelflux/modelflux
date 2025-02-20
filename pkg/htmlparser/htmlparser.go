package htmlparser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/modelflux/modelflux/pkg/util"
	"golang.org/x/net/html"
)

type HTMLParserInputs struct {
	HTML string `yaml:"html"`
}

type HTMLParser struct{}

func (h *HTMLParser) Validate(params map[string]interface{}) error {
	input, err := util.BuildStruct[HTMLParserInputs](params)
	if err != nil {
		return err
	}
	if input.HTML == "" {
		return fmt.Errorf("missing html")
	}
	return nil
}

func (h *HTMLParser) Run(params map[string]interface{}) (string, error) {
	input, err := util.BuildStruct[HTMLParserInputs](params)
	if err != nil {
		return "", err
	}
	extracted, err := extractText(input.HTML)
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}
	return extracted, nil
}

// ExtractText extracts the text content from HTML.
func extractText(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return strings.Join(strings.Fields(buf.String()), " "), nil
}
