package main

import (
	"godb/config"
	"godb/db"
	"fmt"
	"flag"
	"log"
	"godb/utils"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	
)

const (
	SecretKey = "welcome to wangshubo's blog"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

var (
	confPath = flag.String("config", "config.json", "配置文件")
)

func main() {
	flag.Parse()
	conf, err := config.InitConfig(*confPath)
	if err != nil {
		log.Println("load config err:", err)
	}
	fmt.Printf("json string: %s\n",conf.Version)
	fmt.Println(utils.Encrypt("suneil"))
	// db.InitRedisDB()
	db.InitMySqlDB()
	// db.Register("suneil2","abc",0)
	// db.Login("suneil","abc")

	// fmt.Println(time.Now().Unix())
	// fmt.Println(time.Hour * time.Duration(1))
	
	// StartServer()

// 	salt := []byte{0xc8, 0x28, 0xf2, 0x58, 0xa7, 0x6a, 0xad, 0x7b}
// 	dk , err := scrypt.Key([]byte("some password"), []byte(salt), 16384, 8, 1, 32)
// 	fmt.Printf("dk string: %s\n",dk)

// // 	dk, err := scrypt.Key([]byte("some password"), salt, 1<<15, 8, 1, 32)
// // if err != nil {
// //     log.Fatal(err)
// // }
// 	fmt.Println(base64.StdEncoding.EncodeToString(dk))


	// dk, err := scrypt.Key([]byte("some password"), salt, 1<<15, 8, 1, 32)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db.InitRedisDB()
	// db.InitMySqlDB()
	// operationDB()


}

func StartServer() {
	http.HandleFunc("/login", LoginHandler)
	http.Handle("/resource", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))
	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "someone" {
		if user.Password != "p@ssword" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
