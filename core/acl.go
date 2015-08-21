package core

type ACLCheck interface {
	Check(Source string) bool
}
