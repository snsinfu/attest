# attest: Local test tool for competitive coder

[![Release][release-badge]][release-url]
[![Build Status][travis-badge]][travis-url]
[![MIT License][license-badge]][license-url]

attest is a command-line tool that tests a program for given input and expected
output. I made this for quickly checking validity of my program locally when
solving competitive programming problems.

[release-badge]: https://img.shields.io/github/release/snsinfu/attest.svg
[release-url]: https://github.com/snsinfu/attest/releases
[license-badge]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://raw.githubusercontent.com/snsinfu/attest/master/LICENSE.txt
[travis-badge]: https://api.travis-ci.org/snsinfu/attest.svg?branch=master
[travis-url]: https://travis-ci.org/snsinfu/attest

- [Install](#install)
- [Usage](#usage)
- [Options](#options)
- [Testing](#testing)
- [License](#license)


## Install

Download `attest` binary from the [release page][release-url]. Pre-compiled
static executables are listed in the **Assets** box.


## Usage

Create test files in `tests` directory as `*.txt` files. Test file contains
program input and expected output delimited by "---".

```console
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

```console
$ attest ./a.out
PASS  0:01  test1.txt
PASS  0:05  test2.txt
PASS  0:03  test3.txt
```

Each line shows test outcome, execution time and test filename. The test
outcome is one of these:

| Outcome | Meaning                      |
|---------|------------------------------|
| PASS    | Program output was correct   |
| FAIL    | Program output was incorrect |
| TIME    | Program took too long        |
| DEAD    | Program crashed              |

Pass `-v` option to inspect the output of failed tests.

```console
$ attest -v ./a.out
PASS  0:01  test1.txt
FAIL  0:00  test2.txt

FAIL  test2.txt
  Test case
    IN:
    3
    1 2 3
    4 5 6
    7 8 0
    OUT:
    9
  Program output
    OUT:
    0
```


## Options

```
usage: attest [options] <command>...

options:
  -d <tests>    Directory containing test files [default: tests]
  -j <jobs>     Number of concurrent runs; 0 means maximum [default: 0]
  -t <timeout>  Timeout in seconds; 0 means no timeout [default: 0]
  -v            Display detailed information on failed tests
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
