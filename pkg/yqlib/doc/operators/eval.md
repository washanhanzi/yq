# Eval

Use `eval` to dynamically process an expression - for instance from an environment variable.

`eval` takes a single argument, and evaluates that as a `yq` expression. Any valid expression can be used, beit a path `.a.b.c | select(. == "cat")`, or an update `.a.b.c = "gogo"`.

Tip: This can be useful way parameterise complex scripts.

## Dynamically evaluate a path
Given a sample.yml file of:
```yaml
pathExp: .a.b[] | select(.name == "cat")
a:
  b:
    - name: dog
    - name: cat
```
then
```bash
yq 'eval(.pathExp)' sample.yml
```
will output
```yaml
name: cat
```

## Dynamically update a path from an environment variable
The env variable can be any valid yq expression.

Given a sample.yml file of:
```yaml
a:
  b:
    - name: dog
    - name: cat
```
then
```bash
myenv=".a.b[0].name" yq 'eval(strenv(myenv)) = "cow"' sample.yml
```
will output
```yaml
a:
  b:
    - name: cow
    - name: cat
```
