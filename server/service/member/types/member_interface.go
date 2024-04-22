package types

type IMemberSrv interface {
	Auth() (bool, error)
	AuthAndMember() (*Member, error)
	Create(param *MemberCreateParam) error
	Edit(param *MemberEditParam) error
	Delete(param *MemberDeleteParam) error
	Member() (*Member, error)
	Members() ([]Member, error)
}
