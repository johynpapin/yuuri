package enhancer

import (
	"github.com/knakk/rdf"
	"strings"
	"github.com/johynpapin/yuuri/pkg/agrovoc"
	"errors"
	"fmt"
)

// AgrovocEnhancer is an Enhancer capable of linking triples to the Agrovoc database.
type AgrovocEnhancer struct {
	ingredientToAgrovoc map[string]string
	waitingIngredients  map[string][]rdf.Triple
	waitingOthers       map[string][]rdf.Triple
}

// Next takes a triple as parameter and returns all the triples that could be enriched thanks to this triple.
func (ae AgrovocEnhancer) Next(triple rdf.Triple) ([]rdf.Triple, error) {
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
			ae.waitingOthers[triple.Subj.String()] = append(ae.waitingOthers[triple.Subj.String()], triple)
		}
	} else {
		results = append(results, triple)
	}

	return results, nil
}

// End allows the last untreated triples to be returned as they were waiting for a certain triple.
func (ae AgrovocEnhancer) End() ([]rdf.Triple, error) {
	return nil, nil
}

func (ae AgrovocEnhancer) processIngredient(ingredientIRI string) ([]rdf.Triple, error) {
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

		ae.waitingIngredients[ingredientIRI] = []rdf.Triple{}

		return results, nil
	}

	return nil, nil
}

func (ae AgrovocEnhancer) processLabel(label rdf.Triple) ([]rdf.Triple, error) {
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

		return ae.processIngredient(ingredientIRI)
	}

	return []rdf.Triple{}, nil
}

// NewAgrovocEnhancer returns a new AgrovocEnhancer, ready to use.
func NewAgrovocEnhancer() AgrovocEnhancer {
	return AgrovocEnhancer{make(map[string]string), make(map[string][]rdf.Triple), make(map[string][]rdf.Triple)}
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
	return false
}
