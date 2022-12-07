package main

import (
	"database/sql"
	//"log"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id              int
	Nombre          string
	ApellidoPaterno string
	ApellidoMaterno string
	Correo          string
}

// Check here
func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT id, Nombre, ApellidoPaterno, ApellidoMaterno, Correo FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, apellidoPaterno, apellidoMaterno, correo string
		err = registros.Scan(&id, &nombre, &apellidoPaterno, &apellidoMaterno, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.ApellidoPaterno = apellidoPaterno
		empleado.ApellidoMaterno = apellidoMaterno
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Println(arregloEmpleado)

	//fmt.Fprintf(w, "Hello, World!")
	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		apellidoPaterno := r.FormValue("apellidoPaterno")
		apellidoMaterno := r.FormValue("apellidoMaterno")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(Nombre, ApellidoPaterno, ApellidoMaterno, Correo) VALUES (?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, apellidoPaterno, apellidoMaterno, correo)

		http.Redirect(w, r, "/", 301)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT id, Nombre, ApellidoPaterno, ApellidoMaterno, Correo FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, apellidoPaterno, apellidoMaterno, correo string
		err = registro.Scan(&id, &nombre, &apellidoPaterno, &apellidoMaterno, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.ApellidoPaterno = apellidoPaterno
		empleado.ApellidoMaterno = apellidoMaterno
		empleado.Correo = correo
	}
	fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		apellidoPaterno := r.FormValue("apellidoPaterno")
		apellidoMaterno := r.FormValue("apellidoMaterno")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET Nombre=?, ApellidoPaterno=?, ApellidoMaterno=?, Correo=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, apellidoPaterno, apellidoMaterno, correo, id)

		http.Redirect(w, r, "/", 301)
	}
}
