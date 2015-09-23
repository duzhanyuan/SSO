package main

import "fmt"
import "net/http"
import "net/url"
import "io/ioutil"

func doPost(op, id, key string) {
	addr := "http://127.0.0.1:8001/user/"
	if op == "reg" {
		addr = addr + "register"
	} else if op == "login" {
		addr = addr + "login"
	} else {
		addr = addr + "logout"
	}
	resp, err := http.PostForm(addr, url.Values{"key": {key}, "id": {id}})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func main() {
	var op string
	var id string
	var key string
	for {
		fmt.Scanf("%s %s %s", &op, &id, &key)
		doPost(op, id, key)
	}

}
