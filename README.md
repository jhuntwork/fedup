# fedup

```txt
fedup is a simple File dEDUPlicator.

It walks a given directory, collects hash sums of all files,
and turns all duplicates into hard links.

Usage:
  fedup <directory> [flags]

Flags:
  -d, --dryrun   Do not make any changes.
  -h, --help     help for fedup
  -q, --quiet    Do not print any output.
```

`fedup` uses the [BLAKE3](https://github.com/BLAKE3-team/BLAKE3) cryptographic
hash function to compute hash sums, making the operation _fast_.
