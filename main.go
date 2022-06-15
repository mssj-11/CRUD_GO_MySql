package main

// Importando las funciones
import (
	// Formatenado informacion de Entrada/Salida
	"database/sql"  // Conectando a MySql
	"log"           // Informacion en la terminal (hora)
	"net/http"      // Mostrando el sitio
	"text/template" // Separando la informacion en templates

	_ "github.com/go-sql-driver/mysql" // Cargando el driver de MySql github.com/go-sql-driver/mysql
)

// Funcion para la conexion a la DB
func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Password := ""
	NameDB := "sistemago"

	// Condicion de la conexion en caso de algun  error
	conexion, error := sql.Open(Driver, Usuario+":"+Password+"@tcp(127.0.0.1)/"+NameDB) // Remplazar la direccion por la direccion IP en caso sea, o en este caso de manera local es 127.0.0.1

	if error != nil { // SI hay error
		panic(error.Error()) // Mostrando el error
	}
	return conexion
}

// Obtenindo informacion en la carpeta templates y buscar la inf. en esas plantillas
var templates = template.Must(template.ParseGlob("templates/*"))

// Funcion principal
func main() {
	// Accediendo a la funcion Start
	http.HandleFunc("/", Start)          // Escribiendo en el navegador --> localhost, tambien permite entrar a la funcion Start(Inicio)
	http.HandleFunc("/create", Create)   // Cuando acceda al archivo create, se hara uso de la funcion Create
	http.HandleFunc("/insert", Insert)   // Cuando acceda al archivo insert, se hara uso de la funcion Insert
	http.HandleFunc("/edit", Edit)       // Cuando acceda al archivo editar, se hara uso de la funcion Edit
	http.HandleFunc("/update", Update)   // Cuando acceda al archivo actualizar, se hara uso de la funcion Update
	http.HandleFunc("/delete", Delete)   // Cuando acceda al archivo delete, se hara uso de la funcion Delete
	log.Println("Servidor Corriendo...") // Imprimoiendo por consola un mensaje de confirmacion de acceso
	http.ListenAndServe(":8080", nil)    // Iniciando el servidor
}

// Estructura de lectura de la informacion
type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

// Funcion Inicio
/* Agregando los parametros de la funcion, ResponseWriter(Respeta Mayusculas & minusculas), servira para responder a la solicitud
Request, Brindara toda la informacion que estara enviando el usuario
*/
func Start(w http.ResponseWriter, r *http.Request) {

	// Ejecutamdo la conexion a la DB
	conexionEstablacida := conexionDB()
	registros, error := conexionEstablacida.Query("SELECT * FROM empleados")

	if error != nil { // SI hay error
		panic(error.Error()) // Mostrando el error
	}
	// Creando la variable empleado con arreglo de la estructura llamada Empleado
	empleado := Empleado{}
	arregloEmpleado := []Empleado{} // array de emplado sera igual al conjunto de datos de empleados y su estructura

	// Preguntando si se han recorrido los datos
	for registros.Next() {
		var id int
		var nombre, correo string

		// Asignando los valores que se tienen
		error = registros.Scan(&id, &nombre, &correo)

		// Si hay errores
		if error != nil {
			panic(error.Error())
		}
		// Haciendo uso de la estructura
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		// Llenando el arreglo con la informacion
		arregloEmpleado = append(arregloEmpleado, empleado) // Primero asignamos el valor

	}
	// imprimiendo el arreglo
	//fmt.Println(arregloEmpleado)

	// Formateado: (Fprintf)
	//fmt.Fprintf(w, "Hola Bienvenido") // Escribiendo en el navegador, al momento en el que se este accediendo
	templates.ExecuteTemplate(w, "index", arregloEmpleado)
}

//

// Funcion Crear
func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil) // Haciendo uso de la plantilla crear
}

// Funcion Insertar
func Insert(w http.ResponseWriter, r *http.Request) {

	// Condicion, en el caso de que haya un metodo POST
	if r.Method == "POST" {
		// creamos las variables
		nombre := r.FormValue("nombre") // Insertando en el input nombre
		correo := r.FormValue("correo") // Insertando en el input correo

		// Ejecutamdo la conexion a la DB
		conexionEstablacida := conexionDB()
		insertarRegistros, error := conexionEstablacida.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?, ?)")
		// SI hay error
		if error != nil {
			panic(error.Error()) // Mostrando el error
		}

		// Si no hay error entonces se insertaran los registros
		insertarRegistros.Exec(nombre, correo)

		// Redireccionando la informacion (w, r) a la pagina -> "/"
		http.Redirect(w, r, "/", 301)
	}
}

// Funcion editar
func Edit(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") // Recibinedo el valor del metodo get --> id

	conexionEstablacida := conexionDB()
	registro, error := conexionEstablacida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string

		// Asignando los valores que se tienen
		error = registro.Scan(&id, &nombre, &correo)

		// Si hay errores
		if error != nil {
			panic(error.Error())
		}
		// Haciendo uso de la estructura
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	// Mostarndo los datos del registros
	templates.ExecuteTemplate(w, "edit", empleado)

}

// Funcion actualizar
func Update(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// creamos las variables
		id := r.FormValue("id")         // Insertando en el input id
		nombre := r.FormValue("nombre") // Insertando en el input nombre
		correo := r.FormValue("correo") // Insertando en el input correo

		// Ejecutamdo la conexion a la DB
		conexionEstablacida := conexionDB()
		modificarRegistros, error := conexionEstablacida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
		// SI hay error
		if error != nil {
			panic(error.Error()) // Mostrando el error
		}

		// Si no hay error entonces se actualizaran el registro
		modificarRegistros.Exec(nombre, correo, id)

		// Redireccionando la informacion (w, r) a la pagina -> "/"
		http.Redirect(w, r, "/", 301)
	}
}

// Funcion Eliminar
func Delete(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id") // Recibinedo el valor del metodo get --> id
	//fmt.Println(idEmpleado)

	conexionEstablacida := conexionDB()
	borrarRegistro, error := conexionEstablacida.Prepare("DELETE FROM empleados WHERE id=?")
	// SI hay error
	if error != nil {
		panic(error.Error()) // Mostrando el error
	}

	// Si no hay error entonces se eliminara el registro
	borrarRegistro.Exec(idEmpleado)
	// Redireccionando la informacion (w, r) a la pagina -> "/"
	http.Redirect(w, r, "/", 301)
}

/* Creacion del modulo/paquete para la BD
Ingresando el sig. comando: go mod init NAMEMODULO (en mi caso es: system) y se creara el archivo --> go.mod

Agregando el driver para la conexion con MySql y la DB ingresa el si. comando:
go get -u github.com/go-sql-driver/mysql

Se creara el archivo go.sum
*/
