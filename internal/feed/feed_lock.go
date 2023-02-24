package feed

// Lock is a mechanism to run only once the Job that get the feeds and save it in the database.
type Lock struct {
	IsLocked  bool
	Timestamp int64
}
