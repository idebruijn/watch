<img src="https://user-images.githubusercontent.com/141232/54884127-c564e080-4e9f-11e9-9116-3f7b72607beb.gif" width="500" align="right" alt="watch">

# ⏰ watch

_watch_ tool rewritten in go.

## Features

* working aliases
* configurable shell
* windows support

## Usage

```bash
watch [command]
```

## Example
```bash
watch git status
watch 'kubectl get all | grep server'
```

## Install

```bash
go install github.com/idebruijn/watch@latest
```

## Build locally
```bash
go install .
```

Specify command for _watch_ by setting `WATCH_COMMAND` (`bash -cli` by default).

```bash
export WATCH_COMMAND=`fish -c`
```

## License

MIT
