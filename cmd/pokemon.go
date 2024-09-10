package cmd

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/digitalghost-dev/poke-cli/connections"
	"github.com/digitalghost-dev/poke-cli/flags"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var (
	errorColor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#F2055C"))
	errorBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F2055C"))
	helpBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFCC00"))
	styleBold   = lipgloss.NewStyle().Bold(true)
	styleItalic = lipgloss.NewStyle().Italic(true)
)

// ValidateArgs validates the command line arguments
func ValidateArgs(args []string) error {

	if len(args) > 5 {
		return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "\nToo many arguments"))
	}

	if len(args) < 3 {
		return fmt.Errorf(errorBorder.Render(errorColor.Render("Error!"), "\nPlease declare a Pokémon's name after the [pokemon] command", "\nRun 'poke-cli pokemon -h' for more details", "\nerror: insufficient arguments"))
	}

	if len(args) > 3 {
		for _, arg := range args[3:] {
			if arg[0] != '-' {
				errorTitle := errorColor.Render("Error!")
				errorString := fmt.Sprintf("\nInvalid argument '%s'. Only flags are allowed after declaring a Pokémon's name", arg)
				formattedString := errorTitle + errorString
				return fmt.Errorf(errorBorder.Render(formattedString))
			}
		}
	}

	if args[2] == "-h" || args[2] == "--help" {
		flag.Usage()
		os.Exit(0)
	}

	return nil
}

// PokemonCommand processes the Pokémon command
func PokemonCommand() {

	flag.Usage = func() {
		fmt.Println(
			helpBorder.Render(styleBold.Render("USAGE:"), "\n\t", "poke-cli", styleBold.Render("pokemon"), "<pokemon-name>", "[flag]",
				"\n\t", "Get details about a specific Pokémon",
				"\n\t", "----------",
				"\n\t", styleItalic.Render("Examples:"), "\t", "poke-cli pokemon bulbasaur",
				"\n\t\t\t", "poke-cli pokemon flutter-mane --types",
				"\n\t\t\t", "poke-cli pokemon excadrill -t -a",
				"\n",
				styleBold.Render("\nFLAGS:"), "\n\t", "-a, --abilities", "\t", "Prints out the Pokémon's abilities.",
				"\n\t", "-t, --types", "\t\t", "Prints out the Pokémon's typing."),
		)
	}

	flag.Parse()

	pokeFlags, typesFlag, shortTypesFlag, abilitiesFlag, shortAbilitiesFlag := flags.SetupPokemonFlagSet()

	args := os.Args

	err := ValidateArgs(args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	endpoint := strings.ToLower(args[1])
	pokemonName := strings.ToLower(args[2])

	if err := pokeFlags.Parse(args[3:]); err != nil {
		fmt.Printf("error parsing flags: %v\n", err)
		pokeFlags.Usage()
		os.Exit(1)
	}

	_, pokemonName, pokemonID := connections.PokemonApiCall(endpoint, pokemonName, "https://pokeapi.co/api/v2/")
	capitalizedString := cases.Title(language.English).String(pokemonName)

	fmt.Printf("Your selected Pokémon: %s\nNational Pokédex #: %d\n", capitalizedString, pokemonID)

	if *typesFlag || *shortTypesFlag {
		if err := flags.TypesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}

	if *abilitiesFlag || *shortAbilitiesFlag {
		if err := flags.AbilitiesFlag(endpoint, pokemonName); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	}
}