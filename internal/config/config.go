package config

const DefaultExitCommand = "exit"

type Config struct {
	ExitCommand string `json:"exit_command"`
}

func Default() Config {
	return Config{ExitCommand: DefaultExitCommand}
}
