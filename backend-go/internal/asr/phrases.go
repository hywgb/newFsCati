package asr

import (
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Phrases struct {
	Classes map[string][]*regexp.Regexp
}

type rawPhrases struct {
	Version   int                 `yaml:"version"`
	Languages []string            `yaml:"languages"`
	Classes   map[string][]string `yaml:"classes"`
}

func LoadPhrases(path string) (*Phrases, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r rawPhrases
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	p := &Phrases{Classes: map[string][]*regexp.Regexp{}}
	for cls, list := range r.Classes {
		for _, pat := range list {
			re := regexp.MustCompile(pat)
			p.Classes[cls] = append(p.Classes[cls], re)
		}
	}
	return p, nil
}

func (p *Phrases) Match(text string) (string, bool) {
	lower := strings.TrimSpace(text)
	for cls, regs := range p.Classes {
		for _, re := range regs {
			if re.MatchString(lower) {
				return cls, true
			}
		}
	}
	return "", false
}