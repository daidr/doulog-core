package user_attr

type UserAttr struct {
	IsDemoUser bool
}

func ParseUserAttr(attr int) *UserAttr {
	return &UserAttr{
		IsDemoUser: attr&1 == 1,
	}
}
