package poll

type Option struct {
	name  string
	count int
	// List of IDs of users that selected this option
	userIDs []string
}

func (o *Option) AddResult(userIDs []string) {
	o.count = len(userIDs)
	o.userIDs = userIDs
}
