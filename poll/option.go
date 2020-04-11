package poll

type Option struct {
	name  string
	count int
}

func (o *Option) SetCount(i int) {
	o.count = i
}
