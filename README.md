# attest: Local test tool for competitive coder

[![Release][release-badge]][release-url]
[![Build Status][travis-badge]][travis-url]
[![MIT License][license-badge]][license-url]

[release-badge]: https://img.shields.io/github/release/snsinfu/attest.svg
[release-url]: https://github.com/snsinfu/attest/releases
[license-badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://raw.githubusercontent.com/snsinfu/attest/master/LICENSE.txt
[travis-badge]: https://travis-ci.org/snsinfu/attest.svg?branch=master
[travis-url]: https://travis-ci.org/snsinfu/attest

- [Usage](#usage)
- [Options](#options)
- [Testing](#testing)
- [License](#license)


## Usage

Create test files in `tests` directory as `*.txt` files. Test file contains
program input and expected output delimited by "---".

```
$ ls tests
test1.txt  test2.txt  test3.txt

$ cat test1.txt
2
0 1
1 0
---
1
```

Run `attest ./a.out` to test your program `./a.out` against the test files:

```
$ attest ./a.out
PASS test1.txt
PASS test2.txt
PASS test3.txt
```

The test result can be one of these:

| Result | Meaning                      |
|--------|------------------------------|
| PASS   | Program output was correct   |
| FAIL   | Program output was incorrect |
| TIME   | Program took too long        |
| DEAD   | Program crashed              |


## Options

```
usage: attest [options] <command>...

options:
  -d <tests>    Directory containing test files [default: tests]
  -j <jobs>     Number of concurrent runs; 0 means maximum [default: 0]
  -t <timeout>  Timeout in seconds; 0 means no timeout [default: 0]
  -h            Show this message and exit
```


## Testing

Requires Go 1.13.

```
git clone https://github.com/snsinfu/attest
cd attest
make test
```


## License

MIT License
