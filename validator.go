package ProfileBackend

import (
	"fmt"
	pkb "github.com/FarhanRizkiM/pasetobackend"
)

func IsAdmin(Tokenstr, PublicKey string) bool {
	role, err := pkb.DecodeGetRole(PublicKey, Tokenstr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "admin" {
		return false
	}
	return true
}

func IsPK(TokenStr, Publickey string) bool {
	role, err := pkb.DecodeGetRole(Publickey, TokenStr)
	if err != nil {
		fmt.Println("Error : " + err.Error())
	}
	if role != "PK" {
		return false
	}
	return true
}