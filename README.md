# oo

## Example

### show modified files

```
$ oo . "echo {{.}}" 
```

### rsync modified files

```
$ oo -m="nwrc" -s a "rsync -azvc a/{{.Rel}} b/{{.Rel.Dir}}/" 
```

### show deleted files

```
$ oo -m="d" . "echo deleted {{.}}"`
```

## Template

* .     : path of modified file
* .Dir  : directry
* .Rel  : relative path from watched directry
* .Abs  : absolute path
* .Base : file name
* .Ext  : extension

## Build

```
$ ./build.sh
```
