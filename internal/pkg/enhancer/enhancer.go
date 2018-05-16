// Package enhancer contains several Enhancer, i.e. functions to enrich RDF triples.
package enhancer

import "github.com/knakk/rdf"

// Interface Enhancer represents an enhancer, i.e. a system capable of adding information to rdf triples.
type Enhancer interface {
	// Next take a triple, and return new triples. This allows to process the RDF graph as a stream.
	Next(triple rdf.Triple) ([]rdf.Triple, error)

	// End allows to indicate the end of the treatment and returns the last triples that had not yet been processed
	// (because of a lack of information for example).
	End() ([]rdf.Triple, error)

	// Set allows you to adjust the enhancer settings.
	Set(string, interface{}) error

	// Get allows you to read the enhancer settings.
	Get(string) (interface{}, error)
}
