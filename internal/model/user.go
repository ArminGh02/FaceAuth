package model

type User struct {
	Email      string         `bson:"email"`
	Name       string         `bson:"name"`
	IPAddr     string         `bson:"ip_addr"`
	NationalID string         `bson:"national_id"`
	Status     RegisterStatus `bson:"status"`
}

func (u *User) FirstImage() string {
	return u.NationalID + "1"
}

func (u *User) SecondImage() string {
	return u.NationalID + "2"
}
