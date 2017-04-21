# playbulb-candle

An OSX/Linux command line tool for controlling the MiPOW Playbulb Candle. 

Note that the same candle may have a different ID when viewed from different OSX computers.


## Installation

#### Via Homebrew

Installing this was will just get you the CLI, not the examples.

```
brew tap ajesler/playbulb-candle
brew install playbulb-candle
```

#### From source

Install in the usual Go project way. For further information, see https://golang.org/doc/code.html#Command

## Usage

```
# Assuming "e1817cd1d2cd4c088a094b1c31223588" is your candle ID

# Set to a fast rainbow
candle-cli -effect "rainbow" -speed 10 "e1817cd1d2cd4c088a094b1c31223588"

# Set to a solid green
candle-cli -colour "0000FF00" "e1817cd1d2cd4c088a094b1c31223588"

# Set to a pulsing blue
candle-cli -colour "000000FF" -effect "pulse" -speed 120 "e1817cd1d2cd4c088a094b1c31223588"
```

## Examples

#### examples/discover.go

Listens for Bluetooth LE advertising packets and lists any Playbulb Candles found. 
It will keep scanning until it receives a TERM or INT signal.
Useful for finding the ID of your candle which you will need in order to connect to it.

```
$ go run examples/discover.go
Scanning for Playbulb Candles...
Found 'Playbulb Candle 1' with ID '712ef5686bb249b8c12a3ee4def26000'
```

#### examples/demo.go

Shows some of the different modes available on the candle and how to use the library.

## TODO

[ ] Integrate candle discovery into the `candle-cli` command  
[ ] Return a non-zero exit code if connection failed after the timeout  
[ ] Add an 'off' effect that turns the LED off  
[ ] Add a low battery notification  
[ ] Fix the need for the sleep call in `candle-cli/main.go`  
[ ] Improve the usage text  
