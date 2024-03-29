package msg

type Id uint16

const (
	Login  Id = 100
	Logout Id = 101

	CreateRoom Id = 200
	ListRoom   Id = 201
)
