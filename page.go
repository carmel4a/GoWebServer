package main

type IPage interface {
	get() Page
}

type Page struct {
	PageName string
	PageURL  string
}

func (p Page) get() Page {
	return p
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

func getDefaultPage(name string) Page {
	return Page{
		PageName: name,
		PageURL:  "/" + name}
}
