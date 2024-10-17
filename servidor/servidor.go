package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// Criar Usuario insere um usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body) //le o corpo da requisição e retorna a resposta
	if erro != nil {
		w.Write([]byte("Falha na requisição!"))
		return
	}

	var usuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil { //pega o corpo da requisição recebida em json e tenta converter para struct usuario
		w.Write([]byte("Erro ao converter para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados" + erro.Error()))
		return
	}

	stantement, erro := db.Prepare("insert into usuarios (nome, email) values (?, ?)") //prepara uma declaraçaõ sql para ser executada recebendo uma string com placeholders ? para os valores que serão inseridos na tabela
	if erro != nil {
		w.Write([]byte("Erro ao criar o stantement"))
		return
	}

	defer stantement.Close() // garante que a declaração seja fechada após o uso

	insercao, erro := stantement.Exec(usuario.Nome, usuario.Email) //recebe os argumentos e substitui no lugar dos placeholders e executa a declaração no banco
	if erro != nil {
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	idInserido, erro := insercao.LastInsertId() //recupera o ultimo id inserido e armazena na variável idinserido
	if erro != nil {
		w.Write([]byte("Erro ao obter o id inserido"))
		return
	}
	//STATUS CODES
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! Id: %d", idInserido))) // retorna resposta com id inserido
}

// Busca os usuarios existentes no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados" + erro.Error()))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("Select * from usuarios") //executa uma consulta sql e retorna um conjunto de linhas
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario
	for linhas.Next() { //percorre todas as linhas retornadas pela consulta
		var usuario usuario

		if erro := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil { //extrai os dados de cada linha e armazena nos ponteiros referenciados na struct
			w.Write([]byte("Erro ao escanear usuário"))
			return
		}
		usuarios = append(usuarios, usuario) //adiciona o usuario ao slice de usuario

	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarios); erro != nil { //cria um novo encode de json e
		w.Write([]byte("Erro ao converter usuario em JSON"))
		return
	}

}

// Busca usuario específico no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r) //recupera os parametros passados na requisição

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32) //pega o parametro da requisição e converte para inteiro
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados" + erro.Error()))
		return
	}

	linha, erro := db.Query("select * from usuarios where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}

	var usuario usuario
	if linha.Next() {
		if erro := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); erro != nil {
			w.Write([]byte("Erro ao escanear o usuário"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuario); erro != nil {
		w.Write([]byte("Erro ao tentar converter usuario para JSON!"))
		return
	}
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32) //pega o parametro da requisição e converte para inteiro
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	corpoRequisicao, erro := io.ReadAll(r.Body) //le o corpo da requisição e retorna a resposta
	if erro != nil {
		w.Write([]byte("Falha na requisição!"))
		return
	}

	var novoUsuario usuario

	if erro = json.Unmarshal(corpoRequisicao, &novoUsuario); erro != nil { //pega o corpo da requisição recebida em json e tenta converter para struct usuario
		w.Write([]byte("Erro ao converter para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados" + erro.Error()))
		return
	}

	defer db.Close()

	stantement, erro := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o stantement"))
		return
	}

	defer stantement.Close() // garante que a declaração seja fechada após o uso

	if _, erro := stantement.Exec(novoUsuario.Nome, novoUsuario.Email, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	linha, erro := db.Query("select id, nome, email from usuarios where id = ?", ID)
	if erro != nil {
		http.Error(w, "Erro ao buscar usuario atualizado", http.StatusInternalServerError)
		return
	}
	defer linha.Close()

	var usuarioAtualizado usuario
	if linha.Next() {
		if erro := linha.Scan(&usuarioAtualizado.ID, &usuarioAtualizado.Nome, &usuarioAtualizado.Email); erro != nil {
			http.Error(w, "Erro ao escanear o usuário atualizado", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(usuarioAtualizado); erro != nil {
		w.Write([]byte("Erro ao tentar converter usuario atualizado para JSON!"))
		return
	}
}

func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32) //pega o parametro da requisição e converte para inteiro
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados" + erro.Error()))
		return
	}

	defer db.Close()

	stantement, erro := db.Prepare("Delete from usuarios where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao buscar usuario"))
		return
	}

	defer stantement.Close()

	if _, erro := stantement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o usuario"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
