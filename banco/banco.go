package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver de conexão do Sql
)

// Conectar com o banco de dados, drive de conexão com o banco
func Conectar() (*sql.DB, error) {
	stringConexao := "golang:golang@tcp(127.0.0.1:3306)/devbook?charset=utf8&parseTime=true&loc=Local" //usuario, senha, ip e porta, nome do banco,tipo de linguagem local e horario local

	db, erro := sql.Open("mysql", stringConexao) //abre a conexão paramentros drive do banco e string de conexão
	if erro != nil {
		return nil, erro //se o erro for diferente de nulo a conexão falhou então retornamos nulo para db e o erro.
	}

	if erro = db.Ping(); erro != nil { //verifica se conectou
		return nil, erro
	}

	return db, nil
}
