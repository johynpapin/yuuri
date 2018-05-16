package rdfgraph

import "fmt"

// Term represents a Resource, a Literal or a BlankNode
type Term interface {
	String() string
	Equal(Term) bool
}

// Resource represents an URI or an IRI
type Resource struct {
	URI string
}

// String returns the N-Triples representation of the Resource
func (resource Resource) String() string {
	return fmt.Sprintf("<%s>", resource.URI)
}

// Equal returns true if the two Resource are equal
func (resource Resource) Equal(term Term) bool {
	resource2, ok := term.(Resource)
	if !ok {
		return false
	}

	if resource.URI != resource2.URI {
		return false
	}

	return true
}

func (resource Resource) allowedAsASubject()   {}
func (resource Resource) allowedAsAPredicate() {}
func (resource Resource) allowedAsAnObject()   {}

// NewResource returns a new Resource
func NewResource(uri string) *Resource {
	return &Resource{uri}
}

// Literal represents a text value
type Literal struct {
	LexicalForm string
	LanguageTag string
	DataType    Resource // TODO: think about using a term or a resource here (or maybe a pointer on a resource ?)
}

// String returns the N-Triples representation of the Literal
func (literal Literal) String() string {
	// TODO: escape the lexical form
	value := fmt.Sprintf("\"%s\"", literal.LexicalForm)

	if literal.LanguageTag != "" {
		value += "@" + literal.LanguageTag
	} else if !literal.DataType.Equal(Resource{}) {
		value += "^^" + literal.DataType.String()
	}

	return value
}

// Equal returns true if the two Literal are equal
func (literal Literal) Equal(term Term) bool {
	literal2, ok := term.(Literal)
	if !ok {
		return false
	}

	if literal.LexicalForm != literal2.LexicalForm || literal.DataType.URI != literal2.DataType.URI || literal.LanguageTag != literal2.LanguageTag {
		return false
	}

	return true
}

func (literal Literal) allowedAsAnObject() {}

// NewLiteral returns a new Literal
func NewLiteral(lexicalform string) *Literal {
	return &Literal{LexicalForm: lexicalform}
}

// NewLiteralWithLanguageTag returns a new Literal with a LanguageTag
func NewLiteralWithLanguageTag(lexicalform string, languagetag string) *Literal {
	return &Literal{LexicalForm: lexicalform, LanguageTag: languagetag}
}

// NewLiteralWithDataType returns a new Literal with a DataType
func NewLiteralWithDataType(lexicalform string, datatype Resource) *Literal {
	return &Literal{LexicalForm: lexicalform, DataType: datatype}
}

// BlankNode represents an RDF blank node
type BlankNode struct {
	ID string
}

// String returns the N-Triples representation of the BlankNode
func (blankNode BlankNode) String() string {
	return fmt.Sprintf("_:%s", blankNode.ID)
}

// Equal returns true if the two BlankNode are equal
func (blankNode BlankNode) Equal(term Term) bool {
	blankNode2, ok := term.(BlankNode)
	if !ok {
		return false
	}

	if blankNode.ID != blankNode2.ID {
		return false
	}

	return true
}

func (blankNode BlankNode) allowedAsASubject() {}
func (blankNode BlankNode) allowedAsAnObject() {}

// NewBlankNode returns a new BlankNode
func NewBlankNode(id string) *BlankNode {
	return &BlankNode{id}
}

// Subject represents a RDF subject
type Subject interface {
	Term
	allowedAsASubject()
}

// Predicate represents a RDF predicate
type Predicate interface {
	Term
	allowedAsAPredicate()
}

// Object represents a RDF object
type Object interface {
	Term
	allowedAsAnObject()
}

// Triple represents a RDF triple
type Triple struct {
	Subject   Subject
	Predicate Predicate
	Object    Object
}

// NewTriple returns a new Triple
func NewTriple(subject Subject, predicate Predicate, object Object) *Triple {
	return &Triple{subject, predicate, object}
}
