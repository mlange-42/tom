package config

type TerminalColor struct {
	Rune rune
	Tag  string
}

var Colors = []TerminalColor{
	{Rune: ' ', Tag: "[-]"},
	{Rune: 'y', Tag: "[yellow]"},
	{Rune: 'r', Tag: "[red]"},
	{Rune: 'b', Tag: "[blue]"},
	{Rune: 'c', Tag: "[cyan]"},
	{Rune: 'g', Tag: "[green]"},
}
