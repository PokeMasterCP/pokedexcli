# Pokedex CLI

An interactive command-line Pokédex written in Go. It reads from
[PokeAPI](https://pokeapi.co), lets you walk through the game's location areas,
see which Pokémon live in each, try to catch them, and inspect the ones you've
caught.

Built as part of the [boot.dev](https://boot.dev) Go course.

## Requirements

Go 1.26+. No third-party dependencies — everything is standard library.

## Running

```sh
go run .
```

Or build a binary:

```sh
go build -o pokedexcli . && ./pokedexcli
```

## Commands

| Command | Description |
| --- | --- |
| `help` | List the available commands |
| `map` | Show the next 20 location areas |
| `mapb` | Show the previous 20 location areas |
| `explore <area>` | List the Pokémon found in an area |
| `catch <pokemon>` | Try to catch a Pokémon |
| `inspect <pokemon>` | Show height, weight, stats, and types of one you've caught |
| `pokedex` | List everything you've caught |
| `exit` | Quit |

## Example session

```
Pokedex >map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex >explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
- tentacool
- magikarp
- gyarados

Pokedex >catch magikarp
Throwing a Pokeball at magikarp...
magikarp was caught!

Pokedex >inspect magikarp
Name: magikarp
Height: 9
Weight: 100
Stats:
  -hp: 20
  -attack: 10
  ...
Types:
  -water
```

## How it works

**The REPL** reads a line, lowercases and splits it into a command plus an
optional argument, then looks the command up in a `map[string]cliCmd`. Each
command is a struct holding a name, a description, and a callback with the
signature `func(*config, string) error`, so adding a command means adding one
map entry rather than another branch in a switch.

**Paging** works because PokeAPI returns `next` and `previous` URLs with each
page of results. Those are stored as `*string` on the shared `config` — a nil
pointer is a natural "there is no previous page", which is why they aren't plain
strings.

**Catching** is a probability check against the Pokémon's `base_experience`: the
catch succeeds if that value is below a random number in `[0, 650)`, so a
Caterpie is easy and a Dragonite is not. Caught Pokémon live in a map on the
config, which is what `inspect` and `pokedex` read from.

**Caching** sits in `internal/pokecache`. Every API response is stored as raw
JSON keyed by URL, guarded by a mutex, and a background goroutine started by
`NewCache` reaps entries older than the cache interval. Paging back and forth
with `map` and `mapb` therefore only hits the network once per page. 404s are
cached too, using a sentinel value so a repeated typo doesn't cause a fresh
request.

## Project layout

```
main.go              entry point
repl.go              the read-eval-print loop and input cleaning
commands.go          command registry, config, and every command's callback
internal/pokeapi/    PokeAPI HTTP calls and response types
internal/pokecache/  concurrency-safe TTL cache with a reaper goroutine
```

## Tests

```sh
go test ./...
```

Covers input cleaning and the cache's add/get and expiry behaviour.
