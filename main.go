package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	argsWithoutProg := os.Args[1:]
	filePath := argsWithoutProg[0]
	encryptOrDe := argsWithoutProg[1]

	pass := getPasswordCLI("Enter Password: ")
	actual := passToByteKey(pass)

	if encryptOrDe == "e" || encryptOrDe == "-e" ||
		encryptOrDe == "encrypt" || encryptOrDe == "--encrypt" {
		if confirmPassword(pass) {
			encryptFile(filePath, actual)
		} else {
			log.Fatal("Password inputs did not match")
		}
	} else if encryptOrDe == "d" || encryptOrDe == "-d" ||
		encryptOrDe == "decrypt" || encryptOrDe == "--decrypt" {
		if len(argsWithoutProg) > 2 && argsWithoutProg[2] == "f" {
			decryptFile(filePath, actual, true)
		} else {
			decryptFile(filePath, actual, false)
		}
	} else {
		log.Fatal("Encrypt or decrypt option was incorrect")
	}
}

func confirmPassword(first string) bool {
	second := getPasswordCLI("Confirm Password: ")
	return second == first
}

func getPasswordCLI(cliMsg string) string {
	fmt.Println(cliMsg)
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	return string(bytepw)
}

// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
// 32 bytes (AES-256)
func passToByteKey(password string) []byte {
	hashIt := sha256.Sum256([]byte(password))
	return hashIt[0:16]
}

func encryptFile(filepath string, key []byte) {
	infile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	outfile, err := os.OpenFile(filepath+"enc", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	// Append the IV
	outfile.Write(iv)
}

func decryptFile(filepath string, key []byte, createDecFile bool) {
	infile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	fi, err := infile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Fatal(err)
	}

	var outfile *os.File
	if createDecFile {
		handleName := handleOutFileName(filepath)
		outfile, err = os.OpenFile(handleName, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			if createDecFile {
				outfile.Write(buf[:n])
			} else {
				fmt.Println(string(buf[:n]))
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
}

func handleOutFileName(filepath string) string {
	if strings.HasSuffix(filepath, "enc") {
		return filepath[:len(filepath)-3]
	} else {
		return filepath + "dec"
	}
}
