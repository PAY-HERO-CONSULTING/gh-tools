package jobs

type Job interface {
	Body() (string, error) // Body is just a JSON bytes string
	Name() string
	Options() []JobOption
}
