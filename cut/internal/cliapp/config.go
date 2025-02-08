package app

type Flag struct {
	Value string
	Name  string
}

type SelectOption struct {
	OptionType SelectOptionType
	Value      string
}

type SelectOptionType string

const (
	SelectOptionBytes      SelectOptionType = "b"
	SelectOptionCharacters SelectOptionType = "c"
	SelectOptionFields     SelectOptionType = "f"
)

type Config struct {
	Delimeter string
	Select    SelectOption
}
