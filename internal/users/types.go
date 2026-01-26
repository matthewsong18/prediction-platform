package users

type user struct {
	ID          string
	DiscordID   string
	Username    string
	DisplayName string
}

type User interface {
	GetID() string
	GetDiscordID() string
	GetUsername() string
	GetDisplayName() string
}

func (u *user) GetID() string          { return u.ID }
func (u *user) GetDiscordID() string   { return u.DiscordID }
func (u *user) GetUsername() string    { return u.Username }
func (u *user) GetDisplayName() string { return u.DisplayName }

type WinLoss struct {
	Wins   int
	Losses int
}