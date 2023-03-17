# Panacea VC CLI

A easy-to-use CLI for Panacea [vc-sdk](https://github.com/medibloc/vc-sdk)

## Installation

```bash
make install
```

## Examples

```bash
echo '{
  "@context": [
    "https://www.w3.org/2018/credentials/v1"
  ],
  "id": "https://my-verifiable-credential.com",
  "type": "VerifiableCredential"
  "issuer": "did:panacea:76e12ec712ebc6f1c221ebfeb1f",
  "issuanceDate": "0001-01-01T00:00:00Z",
  "credentialSubject": {
    "id": "did:panacea:ebfeb1f712ebc6f1c276e12ec21",
    "first_name": "John",
    "last_name": "Do",
    "nationality": "Korea",
    "age": 21,
    "hobby": "movie"
  },
}' | \
  vccli sign-credential "<panacea-grpc-addr>" "<mnemonic>" | \
  vccli verify-credential "<panacea-grpc-addr>" | \
  vccli sign-presentation "<panacea-grpc-addr>" "<mnemonic>" "<domain>" "<challenge>" | \
  vccli verify-presentation "<panacea-grpc-addr>"
```
