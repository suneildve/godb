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
	http.HandleFunc("/login", cors(LoginHandler))
	http.HandleFunc("/register", cors(RegisterHandler))
	http.Handle("/user",negroni.New(negroni.HandlerFunc(ValidateTokenMiddleware),negroni.Wrap(cors(http.HandlerFunc(ProtectedHandler)))))
	log.Println("Now listening...3001")
	http.ListenAndServe(":3001", nil)
}




func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
		w.Header().Set("content-type", "application/x-www-form-urlencoded")             //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

func SetAllowOrigi(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	// w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	// w.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS, DELETE, PATCH, PUT, HEAD")
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, auth-token, Accept, X-Requested-With") //header的类型
	// w.Header().Set("content-type", "application/x-www-form-urlencoded")             //返回数据格式是json

	w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
	w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
	w.Header().Set("content-type", "application/json;charset=UTF-8")             //返回数据格式是json
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.Header().Set("WWW-Authenticate", `Bearer realm="Dotcoo User Login"`)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println(auth)
	w.Write([]byte(`{"code":0}`))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["account"] = user.Username
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}
	response := Token{tokenString}
	JsonResponse(response, w)
}


func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")  // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
	w.Header().Add("Access-Control-Allow-Credentials", "true") //设置为true，允许ajax异步请求带cookie信息
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") //允许请求方法
	w.Header().Set("content-type", "application/x-www-form-urlencoded")             //返回数据格式是json
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	auth := r.Header.Get("Authorization")
	if auth == "" {
        w.Header().Set("WWW-Authenticate", `Bearer realm="Dotcoo User Login"`)
        w.WriteHeader(http.StatusUnauthorized)
        return
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
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
		return
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
