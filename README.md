# step-tpm-plugin

🔑 `step` plugin for interacting with TPMs. 

> ⚠️ This tool is currently in beta mode and its usage and outputs 
> might change without announcements.

## Simulator

The `step-tpm-plugin` also includes the option to run a TPM simulator.
Running a TPM simulator is useful for development and testing.
It is based on [Microsoft's reference implementation](https://github.com/Microsoft/ms-tpm-20-ref/), and made available through [go-tpm-tools](https://github.com/google/go-tpm-tools/tree/main/simulator).

To enable the TPM simulator, `step-tpm-plugin` first has to be compiled with support enabled:

```console
go build -o step-tpm-plugin -tags tpmsimulator .
```

You can then run the TPM simulator as follows:

```console
$ sudo ./step-tpm-plugin simulator run
+-------------------+------------------------------------------------------------+
| Version           | TPM 2.0                                                    |
| Interface         | command-channel (simulator)                                |
| Manufacturer      | Microsoft (MSFT, 4D534654, 1297303124)                     |
| Vendor Info       | xCG fTPM                                                   |
| Firmware Version  | 8215.1561                                                  |
| EK URI (RSA 2048) | urn:ek:sha256:gBamFwsY+hrTNTQILLrGQ6naornmhl195WZiJo71I7o= |
| UNIX socket       | /var/run/step-tpmsimulator.sock                            |
| Seed              | 3872b33005455e98                                           |
+-------------------+------------------------------------------------------------+
2024/12/05 11:43:21 TPM simulator available at "/var/run/step-tpmsimulator.sock"
```

It requires administrator privileges to create the UNIX socket in `/var/run`. 
You can make it use a different location using `--socket /path/to/tpmsimulator.sock`.

Every time the TPM simulator is started it will act as a new TPM with a new randomly generated RSA EK.
If you need the RSA EK to be stable, you can start the TPM simulator providing a seed previously reported using `--seed <seed>`.

You'll be able to interact with the TPM simulator from other applications supporting TPMs exposed on a UNIX socket.
One example is the `step-tpm-plugin` itself:

```console
$ ./step-tpm-plugin --device ./tpmsimulator.sock info
+------------------+----------------------------------------+
| Version          | TPM 2.0                                |
| Interface        | command-channel                        |
| Manufacturer     | Microsoft (MSFT, 4D534654, 1297303124) |
| Vendor Info      | xCG fTPM                               |
| Firmware Version | 8215.1561                              |
+------------------+----------------------------------------+

./step-tpm-plugin --device ./tpmsimulator.sock ek get
+----------+-------------+------------------------------------------------------------+----------------+
| TYPE     | CERTIFICATE | URI                                                        | CERTIFICATEURL |
+----------+-------------+------------------------------------------------------------+----------------+
| RSA 2048 | -           | urn:ek:sha256:ul3sYf6uQ6jVwtyrQd5XoAuHI10U8gTvEJ6bMj95LXI= |                |
+----------+-------------+------------------------------------------------------------+----------------+
```

### Dynamic Libraries

Compilation of the TPM simulator relies on OpenSSL headers being installed.
You can find some information about this on the [go-tpm-tools](https://github.com/google/go-tpm-tools?tab=readme-ov-file#openssl-errors-when-building-simulator) repository.

It's possible you get an error similar to the one below: 

```console 
dyld[73203]: Library not loaded: @loader_path/../lib/libcrypto.1.1.dylib
  Referenced from: <2E933485-8384-F9D1-D7E5-269A3722A09A> /path/to/step-tpm-plugin
  Reason: tried: '/path/to/step-tpm-plugin/../lib/libcrypto.1.1.dylib' (no such file)
```

In that case you have to adapt the command used to compile the plugin to something like this:

```console
C_INCLUDE_PATH="/path/to/openssl@1.1/1.1.1v/include/" LIBRARY_PATH="/path/to/openssl@1.1/1.1.1v/lib/" go build -o step-tpm-plugin -tags tpmsimulator .
```

#### MacOS

On macOS it's possible to use `brew` to install OpenSSL, and to specify `CGO_CFLAGS` and `CGO_LDFLAGS` for compilation as follows:

```console
brew install openssl
export CGO_CFLAGS="-I$(brew --prefix openssl)/include"
export CGO_LDFLAGS="-L$(brew --prefix openssl)/lib"
```

See [go-tpm-tools](https://github.com/google/go-tpm-tools?tab=readme-ov-file#openssl-errors-when-building-simulator) for additional troubleshooting.

## TODO

Incomplete lists of things still to do:

- Fix / cleanup Makefile
- Tests
- Handle/transform TPM level errors
    - Example: `NCryptCreatePersistedKey returned 8009000F: The operation completed successfully.`
- More functionalities