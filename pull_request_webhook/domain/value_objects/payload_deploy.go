package value_objects

type DeployEvent struct {
	Action  string `json:"action"`
	Repo    Repository `json:"repository"`
	Sender  User `json:"sender"`
	Status  string `json:"status"` // "on" o "off"
	Success bool   `json:"success"` // true si fue exitoso, false si falló
}

/* type Repository struct {
	Name string `json:"name"`
} */

type User struct {
	Login string `json:"login"`
}