# CLI Interpreter Project

This project is a **simple command-line interface (CLI) interpreter** built to explore and learn the features of the Go. The primary goal was to gain hands-on experience with language concepts such as:

- Parsing user input
- Implementing basic interpreter logic
- Managing command execution
- Structuring code for a CLI application

---

## Features

- Basic CLI environment to execute user commands
- Simple input parsing and processing

---

## Purpose

This project served as a learning exercise to:

1. Understand the language's syntax and semantics.
2. Practice handling user input and output effectively.
3. Implement and explore error handling and program control flow.
4. Build a small-scale, functional CLI application for hands-on learning.

---

## How to Build and Run

1. Clone this repository:

   ```bash
   git clone https://github.com/perajarac/cli-interpreter.git
   cd cli-interpreter
   ```

2. Build the project (if applicable):

   ```bash
   go build main.go  # Example for Go
   ```

3. Start typing commands in the CLI.

---

## Usage Example

```bash
$ echo "perajarac"
perajarac

time | tr ":" "." | wc -c >time.txt
8
```

<sup>Explanation: The output stream of the first command time is bound to the input stream of the second command tr, so that all characters
that this command outputs to its output stream is received by the second tr command to its input stream. This one
the second command replaces the ':' character with the '.' character, and the transformed text is received by the third command
wc to your entrance. It again counts all the characters in that text and writes the result to a file time.txt.

```bash
$ version
CLI Interpreter v1.0

```

---

## Project Structure

```
cli-interpreter/
├── main.go            # Entry point of the program
├── reader/            # Module handling user-defined commands and parsing logic
├── file/              # Module handling basic file I/O logic
└── memory/            # Feature tba
└── README.md          # Project documentation
```

---

## Key Learning Outcomes

- Efficient use of the language's features
- Structuring a CLI program in a modular and maintainable way
- Implementing clean input parsing and command execution logic

---

## Future Improvements

- Add support for more advanced commands
- Implement error reporting and recovery
- Allow configuration and extensibility

---



---

## Acknowledgments

Thanks to this project for helping explore language fundamentals in a practical and fun way!
