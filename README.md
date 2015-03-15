# etcvault - proxy for etcd, adding transparent encryption

## Features

- Works as reverse proxy to etcd
  - Can discover other etcd members
  - Support etcd 2.0.x
- Transparent value decryption for GET
- Transparent value encryption for POST, PUT, PATCH
- Multiple keys

## Motivation

Maintaining multiple etcd clusters is hard. We wanted to use same etcd cluster for across services, entire our infrastructure.

But currently etcd has no ACL like feature. All server can read any values even if it's not required for that server (e.g. credentials for different service). That's the reason why I developed Etcvault.

And I know there's ongoing RFC for etcd, about ACL: https://github.com/coreos/etcd/blob/master/Documentation/rfc/api_security.md

## Example

Generate key first.

```
$ mkdir /tmp/keychain
$ etcvault keygen -save /tmp/keychain my-key
```

Start etcd and etcvault.

```
$ etcd -listen-client-urls http://127.0.0.1:2380 &
$ etcvault start -listen http://127.0.0.1:2381 -initial-backends http://127.0.0.1:2379 -keychain /tmp/keychain &
```

Set plain text

```
$ etcdctl --peers http://127.0.0.1:2381 set greeting hello
hello
$ etcdctl get greeting
hello
```

Try encryption/decryption

```
(this means encrypt "hello" with "my-key")
$ etcdctl --peers http://127.0.0.1:2381 set greeting 'ETCVAULT::plain:my-key:hello::ETCVAULT'
hello

$ etcdctl --peers http://127.0.0.1:2381 get greeting
hello

(cannot read directly)
$ etcdctl --peers http://127.0.0.1:2379 get greeting
ETCVAULT::1:my-key::CMOAuEHp/gcbUFvRuQDDMtpIEl/MQ/2OeYT8sluZs8Fc+YjEalDGHzYSn5MM9FafD9fGMHg9ODPYKNk83i1xXZ9zRhKWeuvG8VrU0DlIQ0hdV3px2hDgJppQBYGfr7QVs/0CKaDFUpkMPuhp6dGkzJ+73ZllL3BTb5UjdW3yizYUB82Qs3fwEUZJnLTCvuejxzMF64weInQXnTBkVrt1Mq/QjBWVJvZty8vvAeEHDKo6n5NpgVlZrn48yVHdKWBzO2z5mQO4VK3MPfLUMPQgUsOBqqbUd4N/NjfxCmPL3cO+Y3FD4WiPvbKGGz6IjFnPr7MoWs8etV+vIC/33gOGSQ==::ETCVAULT
```

You can _transform_ `ETCVAULT::...::ETCVAULT` string to proper format using command

```
$ etcvault transform -keychain /tmp/keychain 'ETCVAULT::1:my-key::CMOAuEHp/gcbUFvRuQDDMtpIEl/MQ/2OeYT8sluZs8Fc+YjEalDGHzYSn5MM9FafD9fGMHg9ODPYKNk83i1xXZ9zRhKWeuvG8VrU0DlIQ0hdV3px2hDgJppQBYGfr7QVs/0CKaDFUpkMPuhp6dGkzJ+73ZllL3BTb5UjdW3yizYUB82Qs3fwEUZJnLTCvuejxzMF64weInQXnTBkVrt1Mq/QjBWVJvZty8vvAeEHDKo6n5NpgVlZrn48yVHdKWBzO2z5mQO4VK3MPfLUMPQgUsOBqqbUd4N/NjfxCmPL3cO+Y3FD4WiPvbKGGz6IjFnPr7MoWs8etV+vIC/33gOGSQ==::ETCVAULT'
hello
```

## Detailed Usage

### Generate keys

```
$ etcvault keygen NAME
$ etcvault keygen -save /path/to/keychain/directory NAME
```

for more options, see help.

### Start proxy

```
$ etcvault start -keychain /path/to/keychain/directory -listen http://localhost:2381 -initial-backends http://etcd:2379
```

## Options

- `-listen`: URL to listen to.
- `-advertise-url`: URL to advertise. Used for `/v2/members` and `/v2/machines` response.
- `-keychain`: Path to directory contains key files

### Discovery options

Must be present `-initial-backends` or `-discovery-srv`. Backends are discovered using etcd's API.

- `-initial-backends`: etcd client URLs separated by comma. (e.g. `http://etcd-1:2379,http://etcd-2:2379,...`)
- `-discovery-srv`: FQDN to look up `_etcd-server._tcp` and `_etcd-server-ssl._tcp` SRV records.

### TLS support

etcvault supports HTTPS for both, transport with etcd and listening.

#### Listen https

just specify HTTPS url to `-listen` (e.g. `https://localhost:2381`). Valid certificate options are required.

#### CA and key files

- client:
  - `-client-ca-file`
    - Used to validate etcd client port's server certificate.
    - Also, when etcvault is listening HTTPS, and both `-listen-key-file` `-listen-cert-file` aren't present, this CA certificate will be used to validate etcvault's client certificate.
  - `-client-key-file`, `client-cert-file`
    - Used as client certificate to send to etcd client port.
    - Also, when etcvault is listening HTTPS, and both `-listen-key-file` `-listen-cert-file` aren't present, this certificate will be used as etcvault's server certificate.

- listen:
  - `-listen-ca-file`
    - When present with `-listen-key-file` and `-listen-cert-file`, etcvault will validate its client's certificate using this CA file.
    - (only valid when `-listen-key-file` and `-listen-cert-file` are present)
  - `-listen-key-file`, `listen-cert-file`
    -  When present, etcvault won't use `-client-*` for etcvault's TLS server.
    - This certificate is used for etcvault's server certificate

- peer:
  - `-peer-ca-file`
    - Used to validate etcd peer port's server sertificate.
  - `-peer-key-file`, `peer-cert-file`
    - Used as client certificate to send to etcd peer port.
  - __Note:__ etcvault communicates with etcd peer ports when using `-discovery-srv` option. If you're not using it, you can omit `-peer-*`.

## Key distribution

There's no best way to distribute keys. Try to do with your using server provisioning tools.

Here's what file's required for encryption/decryption:

- Hosts that only encryption
  - Place `${KEYCHAIND_DIR}/${KEY_NAME}.pub`
- Hosts that can do decryption
  - Place `${KEYCHAIND_DIR}/${KEY_NAME}.pem`
  - `${KEY_NAME}.pub` is not necessary.


## FAQ

### Why etcvault communicate with etcd *peer* port?

etcvault communicates with etcd peer port when you're using `-discovery-srv` option. Because SRV records are points to peer port.

## License

MIT License
