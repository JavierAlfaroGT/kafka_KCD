package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

	/*
		El paquete gorilla/muximplementa un enrutador de solicitudes y un despachador para hacer coincidir las solicitudes entrantes con su respectivo controlador
	*/
	"github.com/gorilla/mux"
	/*
		Un puerto Go (golang) del proyecto Ruby dotenv (que carga env vars desde un archivo .env)
	*/
	"github.com/joho/godotenv"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
	}
}

func goGetVariableEnv(keyCode string) string {
	//leo el archivo ".env"
	err := godotenv.Load(".env")
	//si es que hay error al leer el archivo
	if err != nil {
		fmt.Println("No se pudieron leer las variables de entorno :(")
	}
	//devuelvo la variable de entrono solicitada
	return os.Getenv(keyCode)
}

func main() {
	//Mensaje de inicializacion
	fmt.Println("Servidor Kafka iniciado")
	fmt.Println("http://localhost:7000/")

	//Obligando que ese escriba la URL correcta
	router := mux.NewRouter().StrictSlash(true)

	//Definciendo la ruta, URL principal
	router.HandleFunc("/", indexRoute)

	//Definciendo la ruta, URL principal
	router.HandleFunc("/salud", salud)

	//Obtengo el puerto al cual se conecta desde el ".env"
	http_port := ":" + goGetVariableEnv("PORT")

	//Levantamos el server
	//Por si existe un error
	if err := http.ListenAndServe(http_port, router); err != nil {
		fmt.Println(err)
	}
}

func salud(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

/*
Se ejectuta cuando se visita la URL principal
ENDPOINT entre el balanceador y el microservicio
*/
func indexRoute(w http.ResponseWriter, r *http.Request) {
	//Mensaje que aparece al visitar la ruta
	fmt.Fprintf(w, "Bienvenido a Kafka\n")

	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	failOnError(err, "Parsing JSON")
	body["way"] = "Kafka"

	//fmt.Println(body)
	mjson, _ := json.Marshal(body)
	// Publicar el mensaje, convierto el objeto JSON a String
	fmt.Println(string(mjson))
	send_message(string(mjson))
	//Informacion de vuelta que indica que se genero la peticion
	fmt.Fprintf(w, "Mensaje Entregado :)\n")
}

func send_message(datos string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "my-cluster-kafka-bootstrap"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages, beta
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	//for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(datos),
	}, nil)
	//}

	// Wait for message deliveries before shutting down
	//p.Flush(15 * 1000)
}
