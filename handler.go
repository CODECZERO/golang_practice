package main

import("net/http")

func handler(w http.ResponseWriter,r *http.Request){
	responseWithJson(w,200,struct{}{})
}