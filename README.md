## promtool-configmap

Run `promtool` over kubernetes `ConfigMaps` without the fuss of translating between the formats.

- Note: Only `promtool check rules` is supported currently

### Install

```
$ go get github.com/adamdecaf/promtool-configmap
```

### Usage

Pass your `ConfigMap` objects into `promtool-configmap`:

```
$ promtool-configmap testdata/good.yaml
testdata/good.yaml passed checks

# read via stdin
$ cat testdata/bad.yaml | promtool-configmap --
ERROR validating rules for /Users/adam/.../bad.yaml, err=error validating rule file
 Groupname should not be empty
 Group: : unknown fields in rule_group: malformed
```

`promtool-configmap` can also read a json `ConfigMap`:

```
$ promtool-configmap testdata/bad.json
ERROR validating rules for /Users/adam/.../bad.json, err=error validating rule file
 yaml: unmarshal errors:
  line 2: cannot unmarshal !!str `bad` into rulefmt.RuleGroup
```

Note: [`promtool`](https://github.com/prometheus/prometheus/tree/master/cmd/promtool) is required on PATH.
