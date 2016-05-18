package main

import (
    "fmt"
    "io"
    "net/http"
    "os/exec"
    "regexp"
    "runtime"
)

var proxy_servers = []string{
    "10.78.99.8:3728",
    "10.65.13.135:8889",
    "10.142.195.82:8989",
    "10.57.89.177:8889",
    "10.57.91.89:8889",
    "10.167.159.45:8889",
}

var can_used = make(chan string)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/", rootHandler)
    address := "0.0.0.0:8888"
    fmt.Println("Server start on: http://" + address)
    http.ListenAndServe(address, nil)
}

func check(proxy string) { 
    cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("curl -I --connect-timeout 2 --proxy http://%s http://www.baidu.com", proxy))
    result, err := cmd.Output()
    can_use := false
    if err == nil {
        res := fmt.Sprintf("%s", result)
        reg := regexp.MustCompile(`HTTP\/1\.1 200 OK`)
        if reg.MatchString(res) {
            can_use = true
        }
    } 
    if can_use {
        can_used <- proxy
    } else {
        can_used <- ""
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    for i := 0; i < len(proxy_servers); i++ {
        go check(proxy_servers[i])
    }
    io.WriteString(w, "Can used proxy:\n")
    for i := 0; i < len(proxy_servers); i++ {
        res := <- can_used
        if res != "" {
            io.WriteString(w, res + "\n") 
        }
    } 
}
