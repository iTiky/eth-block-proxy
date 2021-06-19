package pkg

// Closer defines a generic Close / Stop interface.
type Closer interface {
	Close() error
}
