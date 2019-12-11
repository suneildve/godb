package db

import (
	"fmt"
	// "time"
)

type User struct {
	Account 		string `db:"account"`
	Password 		string  `db:"password"`
	RegisterTime 	int 	`db:"registerTime"`
	LastLoginTime 	int 	`db:"lastLoginTime"`
	LastLogoutTime 	int 	`db:"lastLogoutTime"`
	BlackTime	 	int 	`db:"blackTime"`
	DeviceId 		string 	`db:"deviceId"`
}

func Register(account, password string, registerTime int) int {
	var list []User
	err := mysqldb.Select(&list,"SELECT account, password, registerTime FROM user WHERE account = ?",account)
	if err != nil {
		fmt.Println(err);
		return 1
	}
	if len(list)>0 {
		fmt.Println("The account already exists");
		return 2
	}
	// result, err := mysqldb.Exec("INSERT INTO user (account, password, registerTime) VALUES (?, ?, ?)",account,password,registerTime)
	_, err = mysqldb.Exec("INSERT INTO user (account, password, registerTime) VALUES (?, ?, ?)",account,password,registerTime)
    if err != nil{
        fmt.Println("insert failed,error： ", err)
        return 3
    }
	fmt.Println("registered successfully")
	return 0
}

func Login(account, password string) User {
	var user User
	err := mysqldb.Get(&user,"SELECT account, password, registerTime FROM user WHERE account = ?",account)
	if err != nil {
		fmt.Println(err);
		return user
	}
	fmt.Printf("%#v\n", user)
	if user.Account == account && user.Password == password {
		return user
	}
	return user
}

