package request

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lucasjones/reggen"
)

var /* const */ variablePattern = regexp.MustCompile("{{.*}}")

type generator interface {
	Generate(limit int) string
}

type Shuffler struct {
	generatorMap map[string]generator
}

func NewShuffler(baseRequest *http.Request) *Shuffler {
	return newShuffler(baseRequest, func(match string) generator {
		gen, err := reggen.NewGenerator(match)
		if err != nil {
			log.Panicf("Invalid RegeExp: %s", match)
		}
		return gen
	})
}

func newShuffler(baseRequest *http.Request, genFunc func(match string) generator) *Shuffler {
	shuffler := &Shuffler{}
	shuffler.generatorMap = make(map[string]generator)
	for _, outerMatch := range variablePattern.FindStringSubmatch(baseRequest.URL.Path) {
		innerMatch := strings.Replace(outerMatch, "{{", "", 1)
		innerMatch = strings.Replace(innerMatch, "}}", "", 1)
		if innerMatch != "" {
			shuffler.generatorMap[outerMatch] = genFunc(outerMatch)
		}
	}

	return shuffler
}

func (s *Shuffler) Shuffle(r *http.Request, params ...int) {
	// NOTE: In a future revision, the parameter won't be needed.
	// This change requires the extraction of a local Request type and wrapping http.Request with it.
	var offset int
	if len(params) > 0 {
		offset = params[0]
	} else {
		offset = -1
	}
	for match, gen := range s.generatorMap {
		switch offset {
		case -1:
			r.URL, _ = r.URL.Parse(strings.Replace(r.URL.Path, match, strings.Trim(gen.Generate(1), "{}"), -1))
		default:
			begin, _ := strconv.Atoi(strings.Trim(match, "{}:"))
			r.URL, _ = r.URL.Parse(strings.Replace(r.URL.Path, match, strconv.Itoa(begin+offset), -1))
		}
	}
}
