package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

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

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]map[string]string, error) {
	var result []map[string]string
	err := yaml.Unmarshal(yml, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildMap(data []map[string]string) map[string]string {
	builtMap := map[string]string {}
	for _, elem := range data {
		builtMap[elem["path"]] = elem["url"]
	}
	return builtMap
}