## promtool-rules-configmap

Run `promtool check rules` on kubernetes `ConfigMaps` without the fuss of translating between the formats.

### Usage

Pass your `ConfigMap` objects into `promtool-rules-configmap`:

```
$ promtool-rules-configmap testdata/good.yaml
testdata/good.yaml passed checks

# read via stdin
$ cat testdata/bad.yaml | promtool-rules-configmap --
ERROR validating rules for /Users/adam/code/src/github.com/adamdecaf/promtool-rules-configmap/testdata/bad.yaml, err=error validating rule file
 Groupname should not be empty
 Group: : unknown fields in rule_group: malformed
```

`promtool-rules-configmap` can also read a json `ConfigMap`:

```
$ promtool-rules-configmap testdata/bad.json
ERROR validating rules for /Users/adam/code/src/github.com/adamdecaf/promtool-rules-configmap/testdata/bad.json, err=error validating rule file
 yaml: unmarshal errors:
  line 2: cannot unmarshal !!str `bad` into rulefmt.RuleGroup
```
