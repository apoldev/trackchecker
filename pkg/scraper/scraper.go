package scraper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apoldev/trackchecker/pkg/scraper/document"
)

const (
	TaskTypeSource        = "source"
	TaskTypeRequest       = "request"
	TaskTypeQuery         = "query"
	DefaultRequestTimeout = 10 * time.Second

	DefaultUserAgent      = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36" //nolint:lll
	DefaultAcceptLanguage = "en-US;q=0.6,en;q=0.4"
	DefaultEncoding       = "utf-8"
)

var (
	ErrDocumentIsNil = errors.New("document is nil")
)

// Scraper can scrape data from delivery service.
type Scraper struct {
	Code  string `json:"code,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

func (s *Scraper) Scrape(args *Args) error {
	start := time.Now()
	defer func() {
		args.ExecuteTime = time.Since(start)
	}()

	for i := range s.Tasks {
		err := s.Tasks[i].Process(args)
		if err != nil {
			// todo analyse err
			return err
		}
	}

	return nil
}

type Task struct {
	Type    string            `json:"type,omitempty"`
	Payload string            `json:"payload,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
	Field   Field             `json:"field,omitempty"`
}

func (t *Task) Process(args *Args) error {
	switch t.Type {
	case TaskTypeSource:
		return t.Source(args)
	case TaskTypeRequest:
		return t.Request(args)
	case TaskTypeQuery:
		return t.Query(args)
	}

	return nil
}

// Source creates document from payload.
func (t *Task) Source(args *Args) error {
	return t.selectDocType([]byte(t.Payload), args)
}

// Request makes http request.
func (t *Task) Request(args *Args) error {
	var err error
	var data []byte
	var method string
	var body io.Reader

	method = "GET"
	if t.Params["method"] != "" {
		method = t.Params["method"]
	}

	if v, ok := t.Params["body"]; ok && method != "GET" {
		replacedBody := args.variables.ReplaceStringFromVariables(v)
		body = strings.NewReader(replacedBody)
	}

	headers := map[string]string{
		"User-Agent":      DefaultUserAgent,
		"Accept-Language": DefaultAcceptLanguage,
		"Encoding":        DefaultEncoding,
	}

	// Если POST но нет content-type - установим дефолтный
	if method == "POST" {
		if _, ok := headers["Content-Type"]; !ok {
			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	link := args.variables.ReplaceStringFromVariables(t.Payload)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultRequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, link, body)

	if err != nil {
		return err
	}

	for i := range headers {
		req.Header.Set(i, args.variables.ReplaceStringFromVariables(headers[i]))
	}

	resp, err := args.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return t.selectDocType(data, args)
}

// Query parses document with xpath or jsonpath or css selector.
func (t *Task) Query(args *Args) error {
	if args.document == nil {
		return ErrDocumentIsNil
	}
	return t.parseDoc(args.document, args.ResultBuilder, &t.Field, t.Field.Path)
}

func (t *Task) parseDoc(doc document.Document, builder *ResultBuilder, field *Field, path string) error {
	var err error
	var node document.Document

	query := field.Query

	if query == "" {
		query = t.Payload
	}

	switch field.Type {
	case FieldTypeArray:
		nodes := doc.FindAll(query)
		for i := range nodes {
			node = nodes[i]
			newPath := path + "." + strconv.Itoa(i)
			if field.Element == nil || field.Element.Type == FieldTypeCommon {
				builder.Set(newPath, node.Value())
			} else if field.Element.Type == FieldTypeObject {
				for j := range field.Element.Object {
					_ = t.parseDoc(node, builder, field.Element.Object[j], newPath+"."+field.Element.Object[j].Path)
				}
			}
		}

	case FieldTypeObject:
		node, err = doc.FindOne(query)
		if err != nil {
			return err
		}
		for j := range field.Object {
			_ = t.parseDoc(node, builder, field.Object[j], path+"."+field.Object[j].Path)
		}
	default:
		node, err = doc.FindOne(query)
		if err != nil {
			return err
		}
		builder.Set(path, node.Value())
	}

	return nil
}

func (t *Task) selectDocType(data []byte, args *Args) error {
	var err error

	switch t.Params["type"] {
	case JSON:
		args.document, err = document.NewJSON(data)
	case HTML:
		args.document, err = document.NewHTML(data)
	case XPATH:
		args.document, err = document.NewHTMLXpath(data)
	case JSONXpath:
		args.document, err = document.NewJSONXpath(data)
	case XML:
		args.document, err = document.NewXMLXpath(data)
	}

	if err != nil {
		return err
	}

	if args.document == nil {
		return ErrDocumentIsNil
	}

	return nil
}
