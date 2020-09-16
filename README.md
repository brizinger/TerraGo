# TerraGo ![CodeScan](https://github.com/brizinger/TerraGo/workflows/Go/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/brizinger/TerraGo)](https://goreportcard.com/report/github.com/brizinger/TerraGo) [![codebeat badge](https://codebeat.co/badges/19552fcd-564d-4be4-a192-1e73fb619792)](https://codebeat.co/projects/github-com-brizinger-terrago-master)

TerraGo is a simple Go Tool that can quickly create a terraform file using user input that is ready for testing. It is useful for testing and playing with new images. 

It currently supports only the Docker provider.


### TODO List

- [ ] Fix Box output in terminal
- [ ] Add support for environmental variables
- [ ] Implement command-like calling with flags

## Installation

Use the run feature of Go to run the tool.

```go run src/TerraGo/main.go```

Or open the already builded file.

```cd bin/```

If on Linux:

```./main``` 

If on Windows: 

```main``` 

## Usage

Run the tool and go through the setup.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPL-3.0](https://choosealicense.com/licenses/gpl-3.0/)
