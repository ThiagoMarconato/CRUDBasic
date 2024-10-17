package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	//CRUD - CREATE, READ, UPDATE, DELETE
	//CREATE - POST
	//READ - GET
	//UPDATE - PUT
	//DELETE - DELETE

	router := mux.NewRouter()                                                               //instancia uma nova rota utilizando importação mux para construir rotas/ HandleFunc defini uma nova rota com url
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)          // definindo rota para inserir usuario, metodo post
	router.HandleFunc("/usuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)         // definindo rota para buscar usuarios, metodo getall
	router.HandleFunc("/usuarios/{id}", servidor.BuscarUsuario).Methods(http.MethodGet)     // definindo rota para buscar usuario, metodo get
	router.HandleFunc("/usuarios/{id}", servidor.AtualizarUsuario).Methods(http.MethodPut)  // definindo rota para dar update em usuario
	router.HandleFunc("/usuarios/{id}", servidor.DeletarUsuario).Methods(http.MethodDelete) // definindo rota para deletar usuario

	fmt.Println("Escutando a porta 5000")
	log.Fatal(http.ListenAndServe(":5000", router)) //listenandserve "escuta" uma porta e uma rota

}
