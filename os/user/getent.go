package user

// NextUserFunc is used in users iteration process. It receives *User for each user record.
// If non-nil error is returned from NextUserFunc - iteration process is terminated.
type NextUserFunc func(*User) error

// NextGroupFunc is used in groups iteration process. It receives *Group for each group record.
// If non-nil error is returned from NextGroupFunc - iteration process is terminated.
type NextGroupFunc func(*Group) error

// IterateUsers iterates over user entries. For each retrieved *User entry provided NextUserFunc is called.
// If CGO is enabled on unix systems IterateUsers is not safe with concurrent usage. Because of this, when
// using IterateUsers with multiple goroutines, locking mechanism such as sync.Mutex must be used in order to
// prevent multiple goroutines calling IterateUsers at the same time.
func IterateUsers(n NextUserFunc) error {
	return iterateUsers(n)
}

// IterateGroups iterates over group entries. For each retrieved *Group entry provided NextGroupFunc is called.
// If CGO is enabled on unix systems IterateGroups is not safe with concurrent usage. Because of this, when
// using IterateGroups with multiple goroutines, locking mechanism such as sync.Mutex must be used in order to
// prevent multiple goroutines calling IterateGroups at the same time.
func IterateGroups(n NextGroupFunc) error {
	return iterateGroups(n)
}
