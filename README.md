# vsort

vsort is sort(1) like command line util for sort version strings.

## Installation

Download from [Github Release](https://github.com/autopp/vsort/releases) or use `go get`:

```
$ go get github.com/autopp/vsort/cmd/vsort
```

## Usage

```
$ vsort -h
Usage:
  vsort [flags] [files]

Flags:
  -h, --help            help for vsort
  -i, --input string    Specify input format. Accepted values are "lines" or "json" (default: "lines"). (default "lines")
  -L, --level int       Expected version level (default -1)
  -o, --output string   Specify output format. Accepted values are "lines" or "json" (default: "lines"). (default "lines")
  -p, --prefix string   Expected prefix of version string.
  -r, --reverse         Sort in reverse order.
      --strict          Make error when invalid version is contained.
  -s, --suffix string   Expected suffix pattern of version string.
  -v, --version         Print the version and silently exits.
```

## Examples

```
$ cat versions.txt
0.1.0
1.0.0
0.10.0
0.2.0

$ vsort versions.txt
0.1.0
0.2.0
0.10.0
1.0.0
```

```
$ vsort < version.txt
0.1.0
0.2.0
0.10.0
1.0.0
```

```
$ echo '["v0.1.0", "v1.0.0", "v0.10.0", "v0.2.0"]' | vsort --input json --output json --prefix v
["v0.1.0","v0.2.0","v0.10.0","v1.0.0"]
```

## License

[Apache License 2.0](LICENSE)

## Author

[@AuToPP](https:/twitter.com/AuToPP)
