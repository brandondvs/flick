package feature

// Flag data representation of a Flag flag
type Flag struct {
	name    string
	enabled bool
}

func (f *Flag) IsEnabled() bool {
	return f.enabled
}

func (f *Flag) Name() string {
	return f.name
}

func (f *Flag) Set(value bool) {
	f.enabled = value
}

func (f *Flag) Toggle() {
	f.enabled = !f.enabled
}

// Create creates a Flag flag with a name and a default value
func Create(name string) *Flag {
	return &Flag{
		name: name,
	}
}
