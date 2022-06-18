package commands

import (
	"github.com/LightningDev1/dgc"
)

type CodeblockCommands struct{}

func (t CodeblockCommands) RegisterCommands(router *dgc.Router) {
	codeBlockCommands := [][]string{
		{"Python", "python", "py"},
		{"Ruby", "ruby", "ruby"},
		{"PHP", "php", "php"},
		{"Perl", "perl", "perl"},
		{"CSS", "css", "css"},
		{"JSON", "json", "json"},
		{"JavaScript", "javascript", "js"},
		{"C++", "cpp", "cpp"},
		{"CSS", "css", "css"},
		{"CSS", "css", "css"},
		{"C", "c", "c"},
		{"C#", "cs", "cs"},
		{"VBScript", "vbcript", "vbscript"},
		{"CoffeeScript", "coffeescript", "coffeescript"},
		{"TypeScript", "typescript", "typescript"},
		{"XML", "xml", "xml"},
		{"YAML", "yaml", "yaml"},
		{"Arduino", "arduino", "arduino"},
		{"CMD", "cmd", "cmd"},
		{"Django", "django", "django"},
		{"DockerFile", "dockerfile", "dockerfile"},
		{"Docker", "docker", "docker"},
		{"Golang", "golang", "golang"},
		{"HTML", "html", "html"},
		{"XHTML", "xhtml", "xhtml"},
		{"CSHTML", "cshtml", "cshtml"},
		{"AutoHotKey", "autohotkey", "autohotkey"},
		{"Bash", "bash", "bash"},
		{"BrainFuck", "brainfuck", "brainfuck"},
		{"CMake", "cmake", "cmake"},
		{"Basic", "basic", "basic"},
		{"Pascal", "pascal", "pascal"},
		{"Elixir", "elixir", "elixir"},
		{"Erlang", "erlang", "erlang"},
		{"Fortran", "fortran", "fortran"},
		{"Godot", "godot", "godot"},
		{"Haskell", "haskell", "haskell"},
		{"TOML", "toml", "toml"},
		{"Kotlin", "kotlin", "kotlin"},
		{"Lua", "lua", "lua"},
		{"SQL", "sql", "sql"},
		{"PostgreSQL", "postgresql", "postgresql"},
		{"Diff", "diff", "diff"},
		{"Fix", "fix", "fix"},
		{"Apache", "apache", "apache"},
		{"INI", "ini", "ini"},
		{"ASCIIDoc", "asciidoc", "asciidoc"},
		{"ARM", "arm", "arm"},
		{"ELM", "elm", "elm"},
	}
	for _, element := range codeBlockCommands {
		router.RegisterCmd(&dgc.Command{
			Name:        element[1],
			Description: "Sends a codeblock with language " + element[0],
			Usage:       "[p]" + element[1] + " <text>",
			Handler:     t.generateFunc(element),
		})
	}
}

func (t CodeblockCommands) generateFunc(element []string) func(ctx *dgc.Ctx) {
	return func(ctx *dgc.Ctx) {
		ctx.RespondText("```" + element[2] + "\n" + ctx.Arguments.Raw() + "\n```")
	}
}
