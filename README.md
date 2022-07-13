# Gaes

Gaes is a command line encryption program that uses Go (golang) and the AES algorithm

# Compiling

Ideally change the binary output name to gaes but you can leave it as the default (main)
```
$ go build -o gaes main.go
```

# Usage

Pass in the file path and then the option to encrypt (e) or decrypt (d). For example:

```
$ ./gaes somefile.txt e
```

# Reference

https://levelup.gitconnected.com/a-short-guide-to-encryption-using-go-da97c928259f