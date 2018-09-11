package urlshort

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

type FormatHandler interface {
	parse() ([]map[string]string, error)
}

type YAMLHandler struct {
	content []byte
}

type JSONHandler struct {
	content []byte
}

func (y *YAMLHandler) parse() ([]map[string]string, error) {
	var result []map[string]string
	err := yaml.Unmarshal(y.content, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (y *JSONHandler) parse() ([]map[string]string, error) {
	var result []map[string]string
	err := json.Unmarshal(y.content, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
} 

func NewFormatHandler(format string, content []byte) FormatHandler {
	switch format {
		case "yaml":
			return &YAMLHandler{content: content}
		case "json":
			return &JSONHandler{content: content}
	}
	return nil
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	mapper := func (w http.ResponseWriter, r *http.Request) {
		url, exists := pathsToUrls[r.URL.Path]
		if exists {
			http.Redirect(w, r, url, 302)
		}
		fallback.ServeHTTP(w, r)
	}
	return http.HandlerFunc(mapper)
}

func GeneralHandler(format string, content []byte, fallback http.Handler) (http.Handler, error) {
	handler := NewFormatHandler(format, content)
	parsedFormat, err := handler.parse()
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedFormat)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(data []map[string]string) map[string]string {
	builtMap := map[string]string {}
	for _, elem := range data {
		builtMap[elem["path"]] = elem["url"]
	}
	return builtMap
}