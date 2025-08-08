# Chronical

A terminal-based puzzle game engine supporting a variety of modes and plain text levelpacks.

Chronical is a terminal-based puzzle game engine that supports a variety of puzzle types.
It uses a plain text, YAML-based format for creating and sharing level packs, making it easy for anyone to create their own puzzles.

## TODO

 - [ ] Convert engine to support update and all gameplay related user input.

 - [ ] Update bindings for games to use keys package?

 - [ ] Update database and log file location to be standard practice and configurable.

 - [ ] Implement standardized configuration

 - [ ] Sudoku engine

 - [ ] Redo menu view to be news-y themed

 - [ ] Swap level validation to engine function. Use that when loading a level to ensure it is valid. Output validation results for each level in a pack when importing.

 - [ ] Worlde Engine

 - [ ] Fix nonogram win validation with known empties. Create test when it does work.


## Installation

To install Chronical, you will need to have Go installed on your system. You can then install Chronical using the following command:

```
go install github.com/nsgnn/chronical@latest
```

## Usage

To start playing, simply run the `chronical` command:

```
chronical
```

### Importing Level Packs

You can import level packs from a YAML file using the `import` command:

```
chronical import /path/to/levelpack.yaml
```

### Exporting Level Packs

You can export your level packs to a YAML file using the `export` command. This is useful for sharing your creations with others:

```
chronical export
```

## Creating Level Packs

Level packs are defined in YAML files. Here is an example of a simple level pack:

```yaml
name: My First Level Pack
author: Nathan
levels:
  - id: 1
    name: Green Hills Zone
    engine: nonogram
    initial: |-
      .........
      .........
      .........
      .........
      .........
      .........
      .........
      .........
      .........
    solution: |-
      .1.1.1.1.
      1.1.1.1.1
      .1.1.1.1.
      1.1.1.1.1
      .1.1.1.1.
      1.1.1.1.1
      .1.1.1.1.
      1.1.1.1.1
      .1.1.1.1.
    width: 9
    height: 9
```

## Development

To get started with development, you will need to have Go installed on your system. You can then clone the repository and install the dependencies:

```
git clone https://github.com/user/chronical.git
cd chronical
```

### Building

To build the project, run the following command:

```
go build -o chronical .
```

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue on GitHub. If you would like to contribute code, please fork the repository and open a pull request.
