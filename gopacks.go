package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"strconv"
)

var html = `<!DOCTYPE html>
<html>
	<head>
		<title>Gopacks.org</title>
	</head>
	<body>
		<h1>Gopacks.org coming soon...</h1>
	</body>
</html>`

var portFlag = flag.Int("port", 80, "The port on which to listen.")

type server struct {
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v requested: %v", r.Host, r.RequestURI)
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	io.WriteString(w, html)
}

func main() {
	flag.Parse()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, os.Kill)

	port := ":" + strconv.Itoa(*portFlag)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	defer listener.Close()
	log.Println("Listening on:", *portFlag)

	log.Println("Starting service...", *portFlag)
	go func() {
		err = fcgi.Serve(listener, &server{})
		if err != nil {
			log.Println("FCGI Serve error:", err)
			exit <- nil
		}
	}()

	sig := <-exit
	if sig != nil {
		log.Println("Received signal:", sig)
	}
	log.Println("Exiting.")
}
