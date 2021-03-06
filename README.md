# arti

`arti` is a simple tool to upload/download/manage artifacts to/from a cloud storage provider
such as S3. It is a more general approach to what [deb-s3](https://github.com/krobertson/deb-s3)
does.
Although `arti` is intended to support multiple storage backends only S3 is supported for now.

`arti` will organize the artifacts according to the following scheme:

```
    <artefact-name>/<artefact-version>/
         <filename>
         <filename>.checksum
    <artefact-name2>/<artefact-version2>/
         <filename2>
         <filename2>.checksum
```

The checksum file will be generated on upload and checked when downloading.
This file contains a algorithm specifier and the hash value seperated by a `:`.


For now the only supported checksum algorithm is SHA-256 which uses the identifier `sha256`
and the hex-encoded hash value:

```
sha256:hex(SHA256(<file>))
```



## Configuration

`arti` uses [cobra](https://github.com/spf13/cobra)/[viper](https://github.com/spf13/viper) for
command-line/config file handling. You may configure `arti` stores using a simple config file
format:

```
stores:
  minio:
    type: "S3"
    endpoint: "127.0.0.1:9000"
    access-key-id: "my-minio-username"
    secret-access-key: "i-will-never-tell-anyone-my-secret"
    location: "us-east-1"
```

Fore more examples see this [file](sample-config.yaml) or [this](sample-config.toml). By default
`arti` looks for the file `$HOME/.arti.(yaml|toml)` but you may override this using the `-c` option.
Besides YAML and TOML you may use JSON and some other syntax. For a full list please read the
[viper documentation](https://github.com/spf13/viper).

Any config option may be overridden using environment variables. Environment variable names start
with the prefix `ARTI_` followed by the key name of the config option where all `.` are replaced
with `__` and all `-` are replaced with `_`.

For example the environment variable for the `secret-access-key` config option of the store
named `minio` looks like this:

```
ARTI_STORES__MINIO__SECRET_ACCESS_KEY=<very-secret>
```



## Examples

All these examples assume that the file `$HOME/.arti.yaml` exists and defines 2 stores named `minio`
and `gcs`.

Store addresses consist of a store name and a bucket within that store. `arti` will create it`s
directory structure within the given bucket. If the bucket does not exist it will be created
when uploading.

### Uploading artefacts

```
arti upload minio/test -n foo -v 1.2.3 foo-1.2.3.tar.gz
```

Output:
```
successfully uploaded 374774 Bytes to 'test:/foo/1.2.3/foo-1.2.3.tar.gz'
sha256:dda436a6ea260e6bf6655688f8f8da34cde6751d4fa720732766868b90858f1d
```

Version numbers must follow the semantic versioning scheme but checks are relaxed. Missing
patch level or even minor number are allowed. The resulting directory will however contain the
full version number:

```
arti put gcs/bar -n foo -v 1.0 my-artefact.zip
```

Output

```
successfully uploaded 374774 Bytes to 'bar:/foo/1.0.0/my-artefact.zip'
sha256:446104b3abefda73b1255e5197702acbeb36b3aada251ad39e1e5ab55292683b
```


### Downloading artefacts

```
arti download minio/test -n foo -v 1.2.3
```

`get` is an alias for the `download` command (as well as `put` is an alias for `upload`)


```
arti get gcs/var -n foo -v 1.0
```


### Listing artefacts

The command `list` or `ls` may be used to list all uploaded artefacts.

```
arti ls minio/test
```

Output:

```
foo 1 366.0kB
hello 2 8.7MB
```

This shows that there is one version of artifact `foo` and 2 versions of artifact `hello`.
The size at the end is the total size of all artifacts of a given name.

You may also list all version of one artifact using the following command

```
arti ls minio/test -n hello
```

Output:

```
1.2.5   8.2MB world.tar.xz
1.2.4 472.1kB world.tar.xz
```

It is also possible to define version ranges:

```
arti ls minio/test -v ">=1.0.0 <2.0.0"
```

of course combining these two filters is also possible

```
arti ls minio/test -n hello -v "1.2.4"
```

Output:

```
1.2.4 472.1kB world.tar.xz
```

For full syntax of version ranges please see this [documentation](https://github.com/blang/semver).


### Deleting artefacts

Deleting artefacts is quite straight-forward:

```
arti del gcs/bar -n foo -v 1
```
