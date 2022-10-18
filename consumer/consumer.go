package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			redisInsert(string(msg.Value))
			mongoInsert(string(msg.Value))

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func redisInsert(datos string) { //ojo va a cambiar

	conn, err := redis.Dial("tcp", "35.223.218.131:6379")
	if err != nil {
		fmt.Println("error de conexión a la base de datos redis", err)
		return
	}
	conn.Send("LPush", "personas", datos)
	rep, err := conn.Do("get", "personas") // Consulta el valor de vaule con la tecla "a" y devuelve
	users, _ := redis.Strings(conn.Do("LRANGE", "personas", 0, -1))
	// Grab one string value and convert it to type byte
	fmt.Println("funciona", rep)
	fmt.Println(users)

}

func mongoInsert(datos string) { //ojo va a cambiar
	responseBody := bytes.NewBuffer([]byte(datos))

	resp, err := http.Post("http://35.184.21.178:3200/", "application/json", responseBody)

	if err != nil {
		log.Println("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	sb := string(body)
	log.Printf(sb)

}
