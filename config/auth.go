package config

type authUser struct {
	Owner    string   //Github user or owner
	APIKey   string   //token or key from secrets
	Projects []string //repos names
}

var Users = []authUser{
	{
		Owner:    "K-Road",
		APIKey:   "secret-github-token",
		Projects: []string{"*"},
	},
}
