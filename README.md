# `encyml` Kubernetes YAML Base64 Encode Secrets Data
Using SOPS and Age with Helm creates an interesting problem when generating or modifying an existing encrypted secrets file.  When defining YAML secrets in a YAML file, they must be base64 encoded.  This prompts developers and managers to manually encode the secrets prior to pasting them into the YAML file.  It can be mistakenly done like this:

```bash
$ echo "secret-value-needing-to-be-encoded" | base64
```
This would output the encoded value AND a line feed within the resultant value.  If this value is copied into the YAML file and used as a secret, it would result in the service or database rejecting the credential or secret due to this additional character.  Line feeds should NOT be introduced here.  The proper command is:

```bash
$ echo -n "secret-value-needing-to-be-encoded" | base64
```
This does not introduce the line feed character.

Obviously, this is an easy mistake to make...

`encyml` can solve this problem and make it easier to encode secrets.

## Installation
`encyml` is a Go program, so it can be easily compiled and installed on any platform.

### Install Go on Mac OS X (if you don't have it already)

```bash
$ brew install go
```

### Install Go on Linux (https://go.dev/dl/)
```
$ curl -OL https://golang.org/dl/go1.21.3.linux-amd64.tar.gz
$ sha256sum go1.21.3.linux-amd64.tar.gz

# Compare the SHA256 checksum to the downloads page
$ sudo tar -C /usr/local -xvf go1.21.3.linux-amd64.tar.gz

# Alter your path in ~/.bashrc or ~/.profile to something like
$ echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
$ source ~/.bashrc

# Check if working
$ go version
> go version go1.16.7 linux/amd64
```

### Install Go on Windows
Visit https://go.dev/dl/ and click on the Windows installer.

### Install `encyml`

```bash
# Clone this repository
$ git clone git@github.com:dwoolworth/encyml.git
$ cd encyml
$ go get gopkg.in/yaml.v2
$ go install encyml.go
```

## Usage For `encyml`
This is a really simple Go app that just removes unnecessary whitespace from the end of each value in the YAML file, and can encode or decode data values for you. It works simplest by reading all the base64 encoded values, decoding them, removing whitespace from the end of each value, and then re-encoding them and writing them back to the same YAML file.  This works as a solution to files in-place.

```bash
$ encyml secrets.yaml
```

It's that simple. You may also want to add additional command line flag:

```bash
$ encyml -o secrets.yaml
```
This outputs the encoded, decoded, and altered values in their original state and trimmed state.  If no whitespace exists, you will see the output unchanged.

For example:

```bash
$ ./encyml -o credentials.yaml
     Key: SPRING_DATASOURCE_VALUE
Original: bXktc2VjcmV0LXBhc3N3b3JkLWZvci1kYg==
 Decoded: my-secret-password-for-db
 Trimmed: my-secret-password-for-db
 Encoded: bXktc2VjcmV0LXBhc3N3b3JkLWZvci1kYg==

     Key: TEST_ENV
Original: dGhpcyBpcyBhIHRlc3QgcGFzc3BocmFzZQ==
 Decoded: this is a test passphrase
 Trimmed: this is a test passphrase
 Encoded: dGhpcyBpcyBhIHRlc3QgcGFzc3BocmFzZQ==

     Key: EXAMPLE_AUTHTOKEN
Original: ZGkyMGQ5dzI5ZWQ5ZjB2dGcwd3JqajQzaDUybGsyMzRoOWR1ZjBzOQ==
 Decoded: di20d9w29ed9f0vtg0wrjj43h52lk234h9duf0s9
 Trimmed: di20d9w29ed9f0vtg0wrjj43h52lk234h9duf0s9
 Encoded: ZGkyMGQ5dzI5ZWQ5ZjB2dGcwd3JqajQzaDUybGsyMzRoOWR1ZjBzOQ==

     Key: AWS_S3_TOKEN
Original: dGVzdF9zM19zZWNyZXQ=
 Decoded: test_s3_secret
 Trimmed: test_s3_secret
 Encoded: dGVzdF9zM19zZWNyZXQ=

File processed successfully!
```
In this case, none of the values were encoded with a line feed.  Contrast with this:

```bash
$ ./encyml -o credentials.yaml
     Key: AWS_S3_TOKEN
Original: dGVzdF9zM19zZWNyZXQK
 Decoded: test_s3_secret

 Trimmed: test_s3_secret
 Encoded: dGVzdF9zM19zZWNyZXQ=

     Key: SPRING_DATASOURCE_VALUE
Original: bXktc2VjcmV0LXBhc3N3b3JkLWZvci1kYg==
 Decoded: my-secret-password-for-db
 Trimmed: my-secret-password-for-db
 Encoded: bXktc2VjcmV0LXBhc3N3b3JkLWZvci1kYg==

     Key: TEST_ENV
Original: dGhpcyBpcyBhIHRlc3QgcGFzc3BocmFzZQ==
 Decoded: this is a test passphrase
 Trimmed: this is a test passphrase
 Encoded: dGhpcyBpcyBhIHRlc3QgcGFzc3BocmFzZQ==

     Key: EXAMPLE_AUTHTOKEN
Original: ZGkyMGQ5dzI5ZWQ5ZjB2dGcwd3JqajQzaDUybGsyMzRoOWR1ZjBzOQ==
 Decoded: di20d9w29ed9f0vtg0wrjj43h52lk234h9duf0s9
 Trimmed: di20d9w29ed9f0vtg0wrjj43h52lk234h9duf0s9
 Encoded: ZGkyMGQ5dzI5ZWQ5ZjB2dGcwd3JqajQzaDUybGsyMzRoOWR1ZjBzOQ==

File processed successfully!
```
The AWS_S3_TOKEN has a line feed after the value because it was encoded incorrectly. After the value is trimmed, you can see it is encoded correctly.

## Use `encyml` To Encode/Decode
You may use `encyml` to encode/decode strings within the YAML file as well.  An example can be tested against the included `decoded.yaml` file:

```bash
$ encyml -e decoded.yaml
```
This will simply base64 encode all values within the yaml file.

```bash
$ encyml -d decoded.yaml
```
This will convert the file back to decoded.

Use the `-o` option flag to view the output.

WARNING:  If you decode a file mutliple times that is already decoded, it will destroy the file.  If a file is encoded multiple times, just decode it multiple times, but decoding is destructive.
