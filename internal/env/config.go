package env

type Config struct {
	Application struct {
		Name  string
		Role  Role
		Stage Stage
	}
	Discovery struct {
		Port int
	}
	HTTP struct {
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

type Role string

const (
	RoleHost Role = "host"
	RolePeer Role = "peer"
)

func (r Role) String() string {
	return string(r)
}
