package main

func main(){
	server:=newApiServer(":3000")
	server.Run()
}