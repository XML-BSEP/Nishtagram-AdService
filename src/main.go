package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
)

func main() {
	//g := gin.Default()
	//
	//g.Run("localhost:8090")


		// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.ProtoVersion = 4
	//cluster.Keyspace="ucstechspecs"
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = session.Query("CREATE KEYSPACE if not exists ucstechspecs WITH replication = { 'class': 'SimpleStrategy', 'replication_factor': '1' };").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	err = session.Query("CREATE TABLE if not exists ucstechspecs.fabricInterconnects (partNumber text, name text, attributes map<text, text>, pictures list<text>, PRIMARY KEY (partNumber));").Exec()
	if err != nil {
		log.Println(err)
		return
	}
	m := make(map[string]string)
	m["color"] = "red"
	m["car"] = "ferrari"

	pictures := []string{"http://fi.1.org", "http://lmgtfy"}

	err = session.Query(`INSERT INTO ucstechspecs.fabricInterconnects (partNumber, name, attributes, pictures) VALUES (?, ?, ?, ?)`,
		"UCSFI-61234", "UCS FI 6120", m, pictures).Exec()

	log.Println("Inserted data")
	if err != nil {
		log.Println("kurcina")
		log.Fatal(err)
	}

	var model string
	iter := session.Query(`SELECT name FROM ucstechspecs.fabricInterconnects WHERE partNumber = ? `,
		"UCSFI-61234").Iter()
	for iter.Scan(&model) {
		fmt.Println("Model: ", model)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

}
