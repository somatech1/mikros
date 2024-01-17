package service

type Env string

const (
	stringEnvNotation = "@env"
)

func NewEnv(name string) Env {
	return Env(name + stringEnvNotation)
}
