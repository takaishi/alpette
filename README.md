# alpette

## server


```
$ cat /etc/alpette/alpette.conf
- foo
  command: echo this is foo
- bar
  command: echo this is bar
```

```
$ alpette server --conf /etc/alpette/alpette.conf
```


## client


```
$ alpette run --task foo
```