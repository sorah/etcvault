# etcvault - proxy for etcd, adding transparent encryption

Still in development phase, so doesn't work. Stay tuned!

## Features

- works as reverse proxy to etcd
- transparent value decryption for GET
- transparent value encryption for POST, PUT
- multiple keys

## Usage

```
$ etcvault -key-dir=/path/to/keys
```

## Value format

### base

- `ETCVAULT::{VERSION}:...::ETCVAULT`
- `ETCVAULT::plain:{KEY}:{PLAINTEXT}::ETCVAULT`

### v1

- `ETCVAULT::1:{KEY}:{FORMAT}:{DATA}::ETCVAULT`
- `ETCVAULT::p1:{KEY}:{PLAIN}::ETCVAULT`

#### Sample

- `ETCVAULT::1:default::AQzrPLEElUzFCnj4Ww5E06CemfAs5rQeWcqI+Ht9aR7JbIwCrgLIzYtlDRxy7qto6WayL7xFh1R+Mw64FOv6JHuhMr121iishiIkPDQnI2foUfHqkRNOjY7bz5p/wF8T4+zgM51EXCBeahWrPUZ4gm8fCJhmOgsmVP9SAQk0zG0=::ETCVAULT`
- `ETCVAULT::1:default:long:AQzrPLEElUzFCnj4Ww5E06CemfAs5rQeWcqI+Ht9aR7JbIwCrgLIzYtlDRxy7qZo6WayL7xFh1R+Mw64FOv6JHuhMr121iishiIkPDQnI2foUfHqkRNOjY7bz5p/wF8T4+zgM51EXCBeahWrPUZ4gm8fCJhmOgsmVP9SAQk0zG0=,AQzrPLEElUzFCnj4Ww5E06CemfAs5rQeWcqI+Ht9aR7JbIwCrgLIzYtlDRxy7qZo6WayL7xFh1R+Mw64FOv6JHuhMr121iishiIkPDQnI2foUfHqkRNOjY7bz5p/wF8T4+zgM51EXCBeahWrPUZ4gm8fCJhmOgsmVP9SAQk0zG0=::ETCVAULT`
- `ETCVAULT::plain:default:p4ssw0rd::ETCVAULT`

## License

Copyright (c) 2015 Shota Fukumori (sora_h)

MIT License

```
Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
