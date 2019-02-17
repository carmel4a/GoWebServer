package main

type Page struct {
	PageName string
	PageURL  string
}

type IPage interface {
	get() Page
}

type LoginPage struct {
	Page
}

func (p LoginPage) get() Page {
	return Page{
		PageName: "login",
		PageURL:  "/login"}
}

type RegisterPage struct {
	Page
}

func (p RegisterPage) get() Page {
	return Page{
		PageName: "register",
		PageURL:  "/register"}
}
