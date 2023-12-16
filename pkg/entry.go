package startup

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type Entry struct {
	path  string
	first []string
	food  map[string][]string
	rest  []string
}

func (e *Entry) String() string {
	sb := &strings.Builder{}
	for _, v := range e.first {
		fmt.Fprintln(sb, v)
	}
	// maps don't preserve order, so annoying
	headers := []string{"早飯", "午飯", "晚飯", "夜宵"}
	for _, section := range headers {
		key := "### " + section
		fmt.Fprintln(sb, key)
		for _, value := range e.food[key] {
			fmt.Fprintln(sb, value)
		}
		fmt.Fprintln(sb)
	}
	for _, v := range e.rest {
		fmt.Fprintln(sb, v)
	}
	return strings.TrimSpace(sb.String()) + "\n"
}

func (e *Entry) FrontMatter() (map[string]any, error) {
	start := false
	fm := []string{}
	for _, v := range e.first {
		border := strings.HasPrefix(v, "---")
		if border {
			if start {
				break
			} else {
				start = true
			}
		} else {
			fm = append(fm, v)
		}
	}
	var v map[string]any
	err := yaml.Unmarshal([]byte(strings.Join(fm, "\n")), &v)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %v", err)
	}
	return v, nil
}

func (e *Entry) WriteToFile() error {
	return os.WriteFile(e.path, []byte(e.String()), 0644)
}

func entryFromFile(p string) (*Entry, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("readFile: %v", err)
	}
	s := strings.Split(string(b), "\n")
	var currentHeader string
	food := make(map[string][]string)
	first := []string{}
	rest := []string{}
	startFood := false
	endFood := false
	what := "我吃了什麼"
	// TODO: wtf is this spaghetti bro
	for _, v := range s {
		if strings.HasPrefix(v, "#") && !strings.HasPrefix(v, "###") {
			hasMatch := strings.Contains(v, what)
			if !startFood && hasMatch {
				startFood = true
			} else if startFood && !hasMatch {
				currentHeader = ""
				endFood = true
			}
		} else if startFood && !endFood && strings.HasPrefix(v, "###") {
			currentHeader = v
			continue
		} else if startFood && !endFood && strings.HasPrefix(v, "-") && currentHeader != "" {
			if _, ok := food[currentHeader]; !ok {
				food[currentHeader] = []string{}
			}
			food[currentHeader] = append(food[currentHeader], v)
			continue
		}

		if v != "" || currentHeader == "" {
			if !endFood {
				first = append(first, v)
			} else if startFood && endFood {
				rest = append(rest, v)
			}
		}
	}
	return &Entry{path: p, first: first, food: food, rest: rest}, nil
}
