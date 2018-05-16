package enhancer

import (
	"errors"
	"fmt"
	"github.com/johynpapin/yuuri/pkg/agrovoc"
	"github.com/knakk/rdf"
	"net/http"
	"os"
	"strings"
)

// AgrovocEnhancer is an Enhancer capable of linking triples to the Agrovoc database.
type AgrovocEnhancer struct {
	ingredientToAgrovoc map[string]string
	waitingIngredients  map[string][]rdf.Triple
	waitingOthers       map[string][]rdf.Triple

	downloadOutput string
}

// Next takes a triple as parameter and returns all the triples that could be enriched thanks to this triple.
func (ae *AgrovocEnhancer) Next(triple rdf.Triple) ([]rdf.Triple, error) {
	var results []rdf.Triple
	var err error

	if objectIsAnIngredient(triple) {
		ingredientIRI := triple.Obj.String()
		ae.waitingIngredients[ingredientIRI] = append(ae.waitingIngredients[ingredientIRI], triple)

		if _, exist := ae.ingredientToAgrovoc[ingredientIRI]; exist {
			results, err = ae.processIngredient(ingredientIRI)
			if err != nil {
				return nil, err
			}
		}
	} else if isPartOfAnIngredient(triple) {
		if isALabel(triple) {
			results, err = ae.processLabel(triple)
			if err != nil {
				return nil, err
			}
		} else if shouldKeep(triple) {
			ingredientIRI := triple.Subj.String()
			ae.waitingOthers[ingredientIRI] = append(ae.waitingOthers[ingredientIRI], triple)

			if _, exist := ae.ingredientToAgrovoc[ingredientIRI]; exist {
				results, err = ae.processOthers(ingredientIRI)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		results = append(results, triple)
	}

	return results, nil
}

// End allows the last untreated triples to be returned as they were waiting for a certain triple.
func (ae *AgrovocEnhancer) End() ([]rdf.Triple, error) {
	return nil, nil
}

// Set allows you to adjust the enhancer settings.
func (ae *AgrovocEnhancer) Set(setting string, value interface{}) error {
	if setting != "download_output" {
		return errors.New("invalid setting")
	}

	valueString, ok := value.(string)
	if !ok {
		return errors.New("invalid value")
	}

	ae.downloadOutput = valueString

	return nil
}

// Get allows you to read the enhancer settings.
func (ae *AgrovocEnhancer) Get(setting string) (interface{}, error) {
	if setting != "download_output" {
		return nil, errors.New("invalid setting")
	}

	return ae.downloadOutput, nil
}

func (ae *AgrovocEnhancer) processIngredient(ingredientIRI string) ([]rdf.Triple, error) {
	if len(ae.waitingIngredients[ingredientIRI]) > 0 {
		rdfAgrovocIRI, err := rdf.NewIRI(ae.ingredientToAgrovoc[ingredientIRI])
		if err != nil {
			return nil, err
		}

		var results []rdf.Triple

		for _, ingredientTriple := range ae.waitingIngredients[ingredientIRI] {
			ingredientTriple.Obj = rdfAgrovocIRI

			results = append(results, ingredientTriple)
		}

		delete(ae.waitingIngredients, ingredientIRI)

		return results, nil
	}

	return nil, nil
}

func (ae *AgrovocEnhancer) processOthers(ingredientIRI string) ([]rdf.Triple, error) {
	if len(ae.waitingOthers[ingredientIRI]) > 0 {
		rdfAgrovocIRI, err := rdf.NewIRI(ae.ingredientToAgrovoc[ingredientIRI])
		if err != nil {
			return nil, err
		}

		var results []rdf.Triple

		for _, otherTriple := range ae.waitingOthers[ingredientIRI] {
			otherTriple.Subj = rdfAgrovocIRI

			results = append(results, otherTriple)
		}

		delete(ae.waitingOthers, ingredientIRI)

		return results, nil
	}

	return nil, nil
}

func (ae *AgrovocEnhancer) processLabel(label rdf.Triple) ([]rdf.Triple, error) {
	obj, ok := (rdf.Term(label.Obj)).(rdf.Literal)
	if !ok {
		return nil, errors.New("invalid label")
	}

	results, err := agrovoc.ExactMatch(obj.String())
	if err != nil {
		return nil, fmt.Errorf("error using agrovoc exactmatch: %s", err)
	}

	ingredientIRI := label.Subj.String()

	if len(results) > 0 {
		ae.ingredientToAgrovoc[ingredientIRI] = results[0].URI

		// Download agrovoc.ttl

		response, err := http.Get(results[0].URI + ".ttl")
		if err != nil {
			return nil, fmt.Errorf("error downloading the agrovoc ttl: %s", err)
		} else {
			defer response.Body.Close()

			downloadOutputFile, err := os.OpenFile(ae.downloadOutput, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				return nil, fmt.Errorf("error openning the agrovoc download output file: %s", err)
			}
			defer downloadOutputFile.Close()

			enc := rdf.NewTripleEncoder(downloadOutputFile, rdf.NTriples)
			dec := rdf.NewTripleDecoder(response.Body, rdf.Turtle)
			triples, err := dec.DecodeAll()
			if err != nil {
				return nil, fmt.Errorf("error decoding the agrovoc turtle: %s", err)
			}
			err = enc.EncodeAll(triples)
			if err != nil {
				return nil, fmt.Errorf("error encoding the agrovoc turtle: %s", err)
			}
		}

		// End

		results1, err := ae.processIngredient(ingredientIRI)
		if err != nil {
			return nil, err
		}

		results2, err := ae.processOthers(ingredientIRI)
		if err != nil {
			return nil, err
		}

		return append(results1, results2...), nil
	}

	return nil, nil
}

// NewAgrovocEnhancer returns a new AgrovocEnhancer, ready to use.
func NewAgrovocEnhancer() AgrovocEnhancer {
	return AgrovocEnhancer{
		ingredientToAgrovoc: make(map[string]string),
		waitingIngredients:  make(map[string][]rdf.Triple),
		waitingOthers:       make(map[string][]rdf.Triple),
	}
}

func objectIsAnIngredient(triple rdf.Triple) bool {
	return triple.Obj.Type() == rdf.TermIRI &&
		strings.HasPrefix(triple.Obj.String(), "https://crwb.sigcidc.org/marmotta/resource/I")
}

func isPartOfAnIngredient(triple rdf.Triple) bool {
	return triple.Subj.Type() == rdf.TermIRI &&
		strings.HasPrefix(triple.Subj.String(), "https://crwb.sigcidc.org/marmotta/resource/I")
}

func isALabel(triple rdf.Triple) bool {
	return triple.Pred.Type() == rdf.TermIRI &&
		triple.Pred.String() == "http://www.w3.org/2000/01/rdf-schema#label"
}

func shouldKeep(triple rdf.Triple) bool {
	return true
}
