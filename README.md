# VoxStream

Low-latency audio streaming proxy server.

## Usage

```
make run            # proxy + mock server + synthetic test stream
make run-file FILE=audio.raw  # stream from PCM file
make build          # build only
make clean          # remove binary
```

Or directly:

```
go run . -stream              # with synthetic stream
go run . -file=audio.raw      # with PCM file
go run .                      # proxy + mock server only
```
