package main

import "github.com/SpayswolfGood/auth/internal/app"

func main() {
	app.LoggerStart()
	go app.Start()
	go app.RunGRPCServer()
	select {}
	//password, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	//if err != nil {
	//	return
	//}
	//fmt.Println(string(password))
}
