package gofiles

func ModelsFileBytes() []byte {
	return []byte(
		`package models

type User struct {
	ID       string ` + "`" + `json:"id,omitempty"` + "`" + `
	User     string ` + "`" + `json:"user" binding:"required"` + "`" + `
	Password string ` + "`" + `json:"password" binding:"required"` + "`" + `
	Books    []Book ` + "`" + `json:"projects" binding:"required"` + "`" + `
}
type Book struct {
	ID     string ` + "`" + `json:"id,omitempty"` + "`" + `
	Type   string ` + "`" + `json:"type"` + "`" + `
	Seller string ` + "`" + `json:"seller" binding:"required"` + "`" + `
	ISBN   string ` + "`" + `json:"isbn" binding:"required"` + "`" + `
}
		`)
}
