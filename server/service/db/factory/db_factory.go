package factory

// import "gorm.io/gorm"

// type gormDbFactory struct{}

// func GetGormDbFactory() *gormDbFactory {
// 	return new(gormDbFactory)
// }

// func (f *gormDbFactory) GetDefaultDbBuilder() IDbBuilder {
// 	return f.GetMysqlDbBuilder()
// }

// func (f *gormDbFactory) GetMysqlDbBuilder() *MysqlDbBuilder {
// 	return nil
// }

// type IDbBuilder interface {
// 	Build() (*gorm.DB, error)
// }

// type MysqlDbBuilder struct{}

// func (b *MysqlDbBuilder) Build() (*gorm.DB, error) { return nil, nil }
