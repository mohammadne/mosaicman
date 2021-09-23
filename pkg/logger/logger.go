package logger

type Logger interface {
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)
}
