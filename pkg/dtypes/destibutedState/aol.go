package destributedstate

type aol struct {
	changes []ichange
}

func newAol() appendOnlyLog {
	a := new(aol)
	a.changes = make([]ichange, 0, 1000)

	return a
}

func (a *aol) append(c ichange) {
	a.changes = append(a.changes, c)
}

func (a *aol) getAll() []ichange {
	return a.changes
}
