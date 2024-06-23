package jwt

type IJWTFactory interface {
	GetJWTAuth() IJWTAuth
}

type defaultJWTFactory struct{}

func newDefaultJWTFactory() IJWTFactory {
	return &defaultJWTFactory{}
}

func (f *defaultJWTFactory) GetJWTAuth() IJWTAuth {
	return newDefaultJWTAuth()
}

func GetJWTFactory() IJWTFactory {
	if true {
		return newDefaultJWTFactory()
	}
	return nil
}
