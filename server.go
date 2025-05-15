package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// Importar el paquete db con el módulo correcto
	"example.com/m/db"
)

// UserRequest es una estructura para recibir los datos del usuario desde HTTP
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func main() {
	// Iniciar conexión con la base de datos
	err := db.Connect()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	// Create a new router
	r := mux.NewRouter()

	// Create a subrouter for API
	apirouter := r.PathPrefix("/api").Subrouter()

	// Ruta para los usuarios
	userRouter := apirouter.PathPrefix("/users").Subrouter()

	// CRUD de usuarios
	userRouter.HandleFunc("", createUserHandler).Methods("POST")
	userRouter.HandleFunc("", getAllUsersHandler).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", getUserHandler).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", updateUserHandler).Methods("PUT")
	userRouter.HandleFunc("/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")

	// Define a route for the profile URL
	r.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Joaquín!, you´ve requested your profile data")
	})

	// Define a route for the profile URL with a variable
	r.HandleFunc("/profile/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		fmt.Fprintf(w, "Hello %s!, you´ve requested your profile data", name)
	})

	// Restrict handler to specific hostname
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Joaquín!, you´ve requested a restricted endpoint")
	}).Host("mydomain.com")

	// Restrict the HTTP methods for the profile URL
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Joaquín!, you´ve requested the test data using %s method", r.Method)
	}).Methods("GET", "POST")

	// Default method for all the other methods
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s HTTP method not allowed", r.Method)
	})

	apirouter.HandleFunc("/characters", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Joaquín!, you´ve requested the characters data")
	}).Methods("GET")

	// Define a route for the static files
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	// Start the server on port 80 and use the router
	log.Println("Servidor iniciado en http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", r))
}

// createUserHandler maneja la creación de un nuevo usuario
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Nombre, email y contraseña son requeridos", http.StatusBadRequest)
		return
	}

	id, err := db.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error al crear el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getAllUsersHandler maneja la obtención de todos los usuarios
func getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, "Error al obtener los usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// getUserHandler maneja la obtención de un usuario específico
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		http.Error(w, "Error al obtener el usuario: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUserHandler maneja la actualización de un usuario
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var user UserRequest
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la petición: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Nombre y email son requeridos", http.StatusBadRequest)
		return
	}

	err = db.UpdateUser(id, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Error al actualizar el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deleteUserHandler maneja la eliminación de un usuario
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = db.DeleteUser(id)
	if err != nil {
		http.Error(w, "Error al eliminar el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
