package main

import (
	"fmt"
	"net/url"
	"os"
)

func init() {
}

// bazel build sps && bazel-bin/sps_/sps "py+mysql://user:123@localhost:3306/database?charset=utf8&tmp=1#test=44&nn=r"
func main() {
	dsn := ""
	if len(os.Args) == 1 {
		help()
		return //exit
	} else {
		dsn = os.Args[1]
	}
	//dsn := "py+mysql://user:123@localhost:3306/database?charset=utf8&tmp=1#test=44&nn=r"
	u, err := parseDSN(dsn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", u)
	fmt.Printf("%#v\n", u.User)

	//db, err := sql.Open("mysql", dsn)
	//if err != nil {
	//	panic(err)
	//}
	//defer func(db *sql.DB) {
	//	_ = db.Close()
	//}(db)
}

// parse url to struct
func parseDSN(dsn string) (u *url.URL, err error) {
	u, err = url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	//if u.Scheme != "mysql" {
	//	return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	//}
	return u, nil
}

func help() {
	fmt.Println("print help content")
}
