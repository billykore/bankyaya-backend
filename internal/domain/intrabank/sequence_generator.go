package intrabank

// SequenceGenerator defines an interface for generating unique sequences.
type SequenceGenerator interface {
	// Generate produces the unique sequence as a string
	// and error if the sequence cannot be generated.
	Generate() (string, error)
}
