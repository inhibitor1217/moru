package env

type Config struct {
	Application struct {
		Name  string
		Stage Stage
	}
	Discovery struct {
		Port int
	}
	Log struct {
		Level string
	}
}

type Stage string

const (
	StageDev        Stage = "dev"
	StageProduction Stage = "production"
)

func (s Stage) String() string {
	return string(s)
}
