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

By default, encryption will have no effect on the input file. You can have the input file 
get deleted after encryption by passing in "d"

```
$ ./gaes somefile.txtenc e d
```

By default, decryption will not create a new file. It will just print the content. To create
a new file pass in "f"

```
$ ./gaes somefile.txtenc d f
```

# Reference

https://levelup.gitconnected.com/a-short-guide-to-encryption-using-go-da97c928259f