# Safe Steps

- define .env variables
- go mod tidy
- go run .

## Env variables

DBEAVER_CONFIG_PATH

``` bash
${HOME}/Library/DBeaverData/workspace6/General/.dbeaver/
```

BOUNDARY_ADDRESS

- Same as the one used for boundary-cli or boundary desktop

## Unsafe steps

- run `boundary-sidecar` executable
