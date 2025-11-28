package main

import (
	"fmt"
	"time"

	"github.com/theStackguy/usercache_go/src"
)

type Person struct {
	Name      string   `json:"full_name"`
	Age       int      `json:"age"`
	IsStudent bool     `json:"is_student"`
	Courses   []string `json:"courses"`
}

func main() {

	um := src.NewUserManager()
	// um.UserFlush(10 * time.Second)
	// nu := cache.NewUserManager()
	userid, err := um.AddNewUser(1* time.Minute)
	// user1,_ := um.AddNewUser(30 * time.Second)
	// user2,_ := nu.AddNewUser(30 * time.Second)

	if err != nil {
		fmt.Println(err)
	}
	p := Person{
		Name:      "Alice Smith",
		Age:       30,
		IsStudent: false,
		Courses:   []string{"History", "Literature"},
	}
	// a := Person{
	// 	Name:      "Anandhu",
	// 	Age:       30,
	// 	IsStudent: false,
	// 	Courses:   []string{"History", "Maths"},
	// }
	err = um.AddOrUpdateUserCache(userid, "newKey",p,0)
	// _= um.AddOrUpdateUserCache(user1,"TestNewKey",p,3*time.Second)
	// 	_ = nu.AddOrUpdateUserCache(user2, "newKey",p,3*time.Second)

	if err != nil {
		fmt.Println(err)
	}
	//  _ = um.RemoveUserCache(userid,"newKey")

	 user,err := um.ReadUser(userid);
	 if err == nil {
		fmt.Println(user)
	 }
	 fmt.Println(err)
	

	 value, err:= util.GenerateRefreshToken(32)
	 fmt.Println(value)
	// err = um.AddOrUpdateUserCache(userid, "newPerson",a,3*time.Second)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// data ,err := um.ReadDataFromCache(userid,"newKey")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// user,err := um.ReadUser(userid);
	// user12,err := um.ReadUser(user1);
	// fmt.Println(data,err)
	// fmt.Println(user)
	// fmt.Println(user12)
	// time.Sleep(30*time.Second)
    // data ,err = um.ReadDataFromCache(userid,"newKey")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(data)
	// data ,err = nu.ReadDataFromCache(user2,"newKey")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(data)
	// user,err = um.ReadUser(userid);
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(user)


	// time.Sleep(10 * time.Second)
	//  data,err := um.ReadUser(userid);
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(data)
	
}
