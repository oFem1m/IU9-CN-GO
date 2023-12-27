package main

import (
	"flag"
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

func main() {
	var (
		root = flag.String("root", "/home/alexandr/BMSTU_git/IU9-CN-GO/lab6.1/server/server_space", "Root directory to serve")
		user = flag.String("user", "Sasha", "Username for login")
		pass = flag.String("pass", "123", "Password for login")
		port = flag.Int("port", 2121, "Port")
		host = flag.String("host", "", "Host")
	)
	flag.Parse()
	if *root == "" {
		log.Fatalf("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
