package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "./proto"

	"google.golang.org/grpc"
)

const address = "grpc:50051"

func main() {
	http.HandleFunc("/", rootHandler)
	log.Println("Listening on :8080...")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		param1 := r.URL.Query().Get("msg")
		if param1 != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{ "status": "OK", "status_code":"1", "message": "` + param1 + `" }`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK", "status_code":"1", "message": "gRPC | Si llegaste acá, ya sabes que hacer."}`))
		return
	case "POST":
		type Person struct {
			Nombre        string `json:"Nombre"`
			Departamento  string `json:"Departamento"`
			Edad          int    `json:"Edad"`
			FormaContagio string `json:"Forma de contagio"`
			Estado        string `json:"Estado"`
		}
		var p Person
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, `{"status":"FAILED","status_code":"0","message":"El cuerpo del mensaje no tiene el formato correcto."}`, http.StatusBadRequest)
			return
		}
		Nombre := p.Nombre
		Departamento := p.Departamento
		Edad := p.Edad
		FormaContagio := p.FormaContagio
		Estado := p.Estado
		mensaje := `{"Nombre":"` + Nombre + `", "Departamento":"` + Departamento + `", "Edad":` + strconv.Itoa(Edad) + `, "FormaContagio":"` + FormaContagio + `", "Estado":"` + Estado + `"}`

		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			http.Error(w, `{"status":"FAILED","status_code":"0","message":"No se pudo conectar. `+err.Error()+`"}`, http.StatusBadRequest)
			return
		}
		defer conn.Close()
		c := pb.NewDataClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SendData(ctx, &pb.SendDataRequest{Data: mensaje})
		if err != nil {
			http.Error(w, `{"status":"FAILED","status_code":"0","message":"No se enviar el mensaje. `+err.Error()+`"}`, http.StatusBadRequest)
		}
		log.Printf("Greeting: %s", r.GetMessage())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"OK", "status_code":"1", "data": {"Enviando": ` + mensaje + `}}`))
	default:
		http.Error(w, `{"status":"FAILED","status_code":"0","message":"Opps, solamente se soportan los métodos GET y POST."}`, http.StatusBadRequest)
		return
	}
}
