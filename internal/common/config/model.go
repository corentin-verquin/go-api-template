package config

type configType string

const (
	String configType = "string"
	Int    configType = "int"
	Bool   configType = "bool"
)

type Config struct {
	Key          string
	DefaultValue any
	Type         configType
}

func NewString(key, def string) Config    { return Config{key, def, String} }
func NewInt(key string, def int) Config   { return Config{key, def, Int} }
func NewBool(key string, def bool) Config { return Config{key, def, Bool} }
