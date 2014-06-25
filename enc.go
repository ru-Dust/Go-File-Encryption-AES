package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"io"
	"io/ioutil"
)

var key []byte


func rand_str(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func createPrivKey() []byte {
	newkey := []byte(rand_str(32))
	err := ioutil.WriteFile("key", newkey, 0644)
	if err != nil {
		fmt.Printf("Error creating Key file!")
		os.Exit(0)
	}
	return newkey
}

func checkKey() {
	thekey, err := ioutil.ReadFile("key") //Check to see if a key was already created
        if err != nil {
                key = createPrivKey() //If not, create one
        } else {
                key = thekey //If so, set key as the key found in the file
        }
}

func main() {
	checkKey()
	fmt.Printf("Key: %x\n", key)
	//fmt.Printf("Data: %s\n", input)
	encryptFile("input", "encryptedfile")
	decryptFile("encryptedfile", "justanotherbot")

}

func encryptFile(inputfile string, outputfile string) {
	b, err := ioutil.ReadFile(inputfile) //Read the target file
        if err != nil {
                fmt.Printf("Unable to open the input file!\n")
                os.Exit(0)
        }
	ciphertext := encrypt(key, b)
        //fmt.Printf("%x\n", ciphertext)
        err = ioutil.WriteFile(outputfile, ciphertext, 0644)
        if err != nil {
                fmt.Printf("Unable to create encrypted file!\n")
                os.Exit(0)
        }
}


func decryptFile(inputfile string, outputfile string) {
	z, err := ioutil.ReadFile(inputfile)
	result := decrypt(key, z)
	//fmt.Printf("Decrypted: %s\n", result)
	fmt.Printf("Decrypted file was created with file permissions 0777\n")
	err = ioutil.WriteFile(outputfile, result, 0777)
	if err != nil {
		fmt.Printf("Unable to create decrypted file!\n")
		os.Exit(0)
	}
}

func encodeBase64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func decodeBase64(b []byte) []byte {
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		fmt.Printf("Error: Bad Key!\n")
		os.Exit(0)
	}
	return data
}

func encrypt(key, text []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	b := encodeBase64(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
	return ciphertext
}
func decrypt(key, text []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(text) < aes.BlockSize {
		fmt.Printf("Error!\n")
		os.Exit(0)
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return decodeBase64(text)
}
