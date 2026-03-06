package feature

// Flag data representation of a feature flag
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

// Create creates a new feature Flag with the given name and default value (false)
func Create(name string) *Flag {
	return &Flag{
		name: name,
	}
}
