package domain

type VisitorRepository interface {
	SaveOrUpdate(entity VisitorEntity)
	Find(ipAddress string) *VisitorEntity
}
