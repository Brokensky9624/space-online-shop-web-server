package types

type IMemberSrv interface {
	Auth(param *MemberAuthParam) error
	AuthAndMember(param *MemberAuthParam) (*Member, error)
	Create(param MemberCreateParam) error
	Edit(param MemberEditParam) error
	Delete(param MemberDeleteParam) error
	Member(param MemberInfoParam) (*Member, error)
	Members() ([]Member, error)
}
