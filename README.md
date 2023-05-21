# jsmap

`jsmap` is a command-line tool written in Go, designed to extract source maps from JavaScript files hosted on any website. It is beneficial for developers, security researchers, and anyone else interested in reviewing the original source code of minified or obfuscated JavaScript files.

## Features
- Extracts source maps from JavaScript files.
- Supports bulk extraction by reading a list of URLs from a file or STDIN.
- Writes extracted source files to a specified directory or separate directories based on the URLs.
- Ignores URLs that do not end in `.js` or `.js?`.

## Installation

Before you start, make sure you have Go installed on your machine. You can download it from the official [Golang website](https://golang.org/dl/).

### Quick install

```
go get github.com/zerodaywolf/jsmap
```

### Building from Source

Clone the repository to your local machine:

```
https://github.com/zerodaywolf/jsmap.git
```

Then, navigate to the `jsmap` directory and build the project:

```
cd jsmap
go build
```

## Usage

To extract source maps from a list of URLs provided via STDIN:

```
$ echo "https://example.com/script.js" | jsmap
$ cat urls.txt | jsmap
```

To extract source maps from a list of URLs provided in a file:

```
$ jsmap -f urls.txt
```

To specify a directory for the output:

```
$ jsmap -o ./output
```

To specify both a file for input and a directory for output:

```
$ jsmap -f urls.txt -o ./output
```

In all cases, the directories will be named after the URL, with slashes (`/`) replaced by underscores (`_`).

## Contributing

Contributions from the community are welcome. If you would like to contribute, please fork the repository, make your changes, and open a pull request. If you have any questions or need help, please feel free to open an issue.

## License

Jsmap is licensed under the BSD 3-Clause "New" or "Revised" License. See the `LICENSE` file for more details.

## Disclaimer

Please use this tool responsibly, and ensure you have permission to extract source maps from any website you target. We are not responsible for any misuse of this tool.
