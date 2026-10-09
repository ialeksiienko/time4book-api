package authport

type Commands struct {
	Login    Login
	Logout   Logout
	Refresh  Refresh
	Register Register
}
