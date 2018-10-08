# alpette

## server


```
$ cat /etc/alpette/alpette.conf
[[tasks]]
name = "foo"
command = "echo this is foo"

[[tasks]]
name = "bar"
command = "echo this is bar"
```

```
$ alpette server --conf /etc/alpette/alpette.conf
```


## client


```
$ alpette run --task foo
```