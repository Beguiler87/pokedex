// This project was great fun to work on, and is fully functional.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedex/internal/pokecache"
	"strings"
	"time"
)

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\nhelp: Displays a help message\nexit: Exit the Pokedex\npokedex: Displays a list of all the Pokemon you've caught\nmap: Displays the next 20 names of locations in the Pokemon world\nmapb: Displays the last 20 names of locations in the Pokemon world. If you are on the first page, the system will notify you.\nexplore <area-name>: Displays a list of Pokemon found in the indicated area.\ncapture: Attempts to capture the indicated Pokemon\ninspect: Displays information about the selected Pokemon based on the Pokedex entry")
	return nil
}
func commandMap(cfg *Config, args []string) error {
	url := ""
	if cfg.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = cfg.Next
	}
	cachedData, found := cfg.Cache.Get(url)
	var body []byte
	resp := locationAreaResponse{}
	if found {
		body = cachedData
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %v", url, err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, body)
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nework response: %v", err)
		}
	}
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	if resp.Next != nil {
		cfg.Next = *resp.Next
	} else {
		cfg.Next = ""
	}
	if resp.Previous != nil {
		cfg.Previous = *resp.Previous
	} else {
		cfg.Previous = ""
	}
	return nil
}
func commandMapb(cfg *Config, args []string) error {
	url := ""
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = cfg.Previous
	}
	cachedData, found := cfg.Cache.Get(url)
	var body []byte
	resp := locationAreaResponse{}
	if found {
		body = cachedData
		err := json.Unmarshal(body, &resp)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %v", url, err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, body)
		err = json.Unmarshal(body, &resp)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nework response: %v", err)
		}
	}
	for _, area := range resp.Results {
		fmt.Println(area.Name)
	}
	if resp.Previous != nil {
		cfg.Previous = *resp.Previous
	} else {
		cfg.Previous = ""
	}
	if resp.Next != nil {
		cfg.Next = *resp.Next
	} else {
		cfg.Next = ""
	}
	return nil
}
func commandExplore(cfg *Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide a location area name: e.g. 'explore pastoria-city-area'")
		return nil
	}
	areaName := strings.Join(args, "-")
	url := "https://pokeapi.co/api/v2/location-area/" + areaName
	body, found := cfg.Cache.Get(url)
	if !found {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %v", url, err)
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, body)
	}
	var area LocationArea
	err := json.Unmarshal(body, &area)
	if err != nil {
		return fmt.Errorf("failed to unmarshal location area: %v", err)
	}
	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
func commandCatch(cfg *Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide the name of an available Pokemon: e.g. 'catch pikachu'")
		return nil
	}
	pokemonName := strings.Join(args, "-")
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	body, found := cfg.Cache.Get(url)
	if !found {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to GET %s: %v", url, err)
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		cfg.Cache.Add(url, body)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	var pmon Pokemon
	err := json.Unmarshal(body, &pmon)
	if err != nil {
		return fmt.Errorf("failed to unmarshal pokemon info: %v", err)
	}
	chance := 100 - (pmon.BaseExperience / 10)
	if chance < 1 {
		chance = 1
	}
	roll := rand.Intn(100)
	if roll < chance {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.Pokedex[pokemonName] = pmon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}
func commandInspect(cfg *Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("Please provide the name of a Pokemon: e.g. 'inspect pikachu'")
		return nil
	}
	pokemonName := strings.Join(args, "-")
	entry, found := cfg.Pokedex[pokemonName]
	if !found {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %v\nHeight: %v\nWeight: %v\nStats:\n", entry.Name, entry.Height, entry.Weight)
	for _, t := range entry.Stats {
		fmt.Printf("  -%v: %v\n", t.Stat.Name, t.Basestat)
	}
	fmt.Println("Types:")
	for _, t := range entry.Types {
		fmt.Printf("  -%v\n", t.Type.Name)
	}
	return nil
}
func commandPokedex(cfg *Config, args []string) error {
	fmt.Println("Your Pokedex:")
	for pokemonName, _ := range cfg.Pokedex {
		fmt.Printf(" - %v\n", pokemonName)
	}
	return nil
}

type Config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
	Pokedex  map[string]Pokemon
}
type cliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}
type locationAreaResponse struct {
	Results  []locationAreaResult `json:"results"`
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
}
type locationAreaResult struct {
	Name string `json:"name"`
}
type LocationArea struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}
type PokemonEncounter struct {
	Pokemon PokemonInfo `json:"pokemon"`
}
type PokemonInfo struct {
	Name string `json:"name"`
}
type Pokemon struct {
	BaseExperience int        `json:"base_experience"`
	Name           string     `json:"name"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []statJSON `json:"stats"`
	Types          []TypeInfo `json:"types"`
}
type Stat struct {
	Name  string
	Value int
}
type TypeInfo struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}
type statJSON struct {
	Basestat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	config := &Config{}
	config.Pokedex = make(map[string]Pokemon)
	interval := 10 * time.Second
	config.Cache = pokecache.NewCache(interval)
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 names of locations in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the last 20 names of locations in the Pokemnon world. If you are on the first page, the system will notify you.",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays a list of Pokemon found in the chosen area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to apture the indicated Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information about the selected Pokemon based on the Pokedex entry",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of all the Pokemon you've caught",
			callback:    commandPokedex,
		},
	}
	Scanner := bufio.NewScanner(os.Stdin)
	for {
		_ = Scanner.Scan()
		scannedText := Scanner.Text()
		loweredText := strings.ToLower(scannedText)
		trimmed := strings.TrimSpace(loweredText)
		words := strings.Split(trimmed, " ")
		firstWord := words[0]
		command, exists := commands[firstWord]
		if exists {
			err := command.callback(config, words[1:])
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}
