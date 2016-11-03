# arti

`arti` is a simple tool to upload/download/manage artifacts to/from a cloud storage provider
such as S3. It is a more general approach to what [deb-s3](https://github.com/krobertson/deb-s3)
does.
Although `arti` is intended to support multiple storage backends only S3 is supported for now.

`arti` will organize the artifacts according to the following scheme:

```
    <artefact-name>/<artefact-version>/
         <filename>
         sha256sum
    <artefact-name2>/<artefact-version2>/
         <filename2>
         sha256sum
```

The sha256sum will be generated on upload and checked when downloading.

`arti` uses a simple config file which can be used to define so called stores. See
[this example](sample-config.toml) for more details. By default `arti` looks for the file
`$HOME/.arti.toml` but you may override this using the `-c` option. Besides toml you may
use json or yaml as a config file format.

Store addresses consist of a store name and a bucket within that store. `arti` will create it`s
directory structure within the given bucket. If the bucket does not exist it will be created
when uploading.



## Examples

All these examples assume that the file `$HOME/.arti.toml` exists and defines 2 stores named `minio`
and `gcs`.


### Uploading artefacts

```
arti upload minio/test -n foo -v 1.2.3 foo-1.2.3.tar.gz
```

Output:
```
successfully uploaded 374774 Bytes to 'test:/foo/1.2.3/foo-1.2.3.tar.gz'
SHA256: k4odDkqLqMvpPdlU75K3wFWVCOxPoI9AyaL4o8dra+8=
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
SHA256: k4odDkqLqMvpPdlU75K3wFWVCOxPoI9AyaL4o8dra+8=
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
foo:
  1.2.3:  366.0 kB  foo-1.2.3.tar.gz
hello:
  1.2.4:  472.1 kB  world.tar.xz
  1.2.5:    8.2 MB  world.tar.xz
```


### Deleting artefacts

Deleting artefacts is quite straight-forward:

```
arti del gcs/bar -n foo -v 1
```