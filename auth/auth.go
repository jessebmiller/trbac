package auth

type Auth struct {
	Permissions Permissionser
	Constraints ConstraintRunner
}

func (a Auth) May (c Context) bool {
	return false
}


