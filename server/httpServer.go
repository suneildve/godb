package server

import (
	"net/http"
	"godb/config"
	"fmt"
	"log"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"encoding/json"
	// "strings"
	"time"
	// "encoding/base64"

)

const (
	SecretKey = "suneil"
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

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


// type User struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}


// 'msg': 'success',
//       'code': 0,
//       'data': {
//         'token': '4344323121398'
//         // 其他数据
// 	  }
	  


func StartHTTP() {
	var conf = config.GetConfig()
	if post, ok := conf.Server["post"]; ok {
		fmt.Println(post)
	}
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.Handle("/user", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))
	log.Println("Now listening...")
	http.ListenAndServe(":3001", nil)
}

func SetAllowOrigi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS, DELETE, PATCH, PUT, HEAD")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, auth-token, Accept, X-Requested-With") //header的类型
	w.Header().Set("content-type", "application/json;charset=utf-8")             //返回数据格式是json
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	SetAllowOrigi(w,r)
	// json, err := json.Marshal(response)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write(json)
	w.Write([]byte(`{"code":0}`))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
        w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
        w.WriteHeader(http.StatusUnauthorized)
        // return
    }
	fmt.Println(auth)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	SetAllowOrigi(w,r)
	var user UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}
	// if strings.ToLower(user.Username) != "someone" {
	// 	if user.Password != "p@ssword" {
	// 		w.WriteHeader(http.StatusForbidden)
	// 		fmt.Println("Error logging in")
	// 		fmt.Fprint(w, "Invalid credentials")
	// 		return
	// 	}
	// }

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["account"] = user.Username
	token.Claims = claims

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w, "Error extracting the key")
	// 	fatal(err)
	// }

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(json)
	// JsonResponse(response, w)
}
func LoginHandler00(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
        w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
        w.WriteHeader(http.StatusUnauthorized)
        // return
    }
	fmt.Println(auth)
    // auths := strings.SplitN(auth, " ", 2)
    // if len(auths) != 2 {
    //     fmt.Println("error")
    //     return
    // }
	// authMethod := auths[0]
    // authB64 := auths[1]
    // switch authMethod {
    // case "Basic":
    //     authstr, err := base64.StdEncoding.DecodeString(authB64)
    //     if err != nil {
    //         fmt.Println(err)
    //         // io.WriteString(w, "Unauthorized!\n")
    //         return
    //     }
    //     fmt.Println(string(authstr))
    //     userPwd := strings.SplitN(string(authstr), ":", 2)
    //     if len(userPwd) != 2 {
    //         fmt.Println("error")
    //         return
    //     }
    //     username := userPwd[0]
    //     password := userPwd[1]
    //     fmt.Println("Username:", username)
    //     fmt.Println("Password:", password)
	// 	fmt.Println()
	// case "Bearer":
    //     // authstr, err := base64.StdEncoding.DecodeString(authB64)
    // default:
    //     fmt.Println("error")
    //     return
    // }

	var user UserCredentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}
	fmt.Println(user)
		// w.Write(json)
	SetAllowOrigi(w,r)
	
	// fmt.Fprintf(w, )
	w.Write([]byte(`{"code":0}`))
    // fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}


func StartServer() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.Handle("/", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(IndexHandler)),
	))
	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}



func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	SetAllowOrigi(w,r)
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func LoginHandler01(w http.ResponseWriter, r *http.Request) {
	var user UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}
	// if strings.ToLower(user.Username) != "someone" {
	// 	if user.Password != "p@ssword" {
	// 		w.WriteHeader(http.StatusForbidden)
	// 		fmt.Println("Error logging in")
	// 		fmt.Fprint(w, "Invalid credentials")
	// 		return
	// 	}
	// }

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	// claims["name"] = "suneil"
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
	SetAllowOrigi(w,r)
	auth := r.Header.Get("Authorization")
	if auth == "" {
        w.Header().Set("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
        w.WriteHeader(http.StatusUnauthorized)
        // return
    }
	fmt.Println(auth)
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
