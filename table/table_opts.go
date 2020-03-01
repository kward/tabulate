package table

type options struct {
	commentPrefix  string
	enableComments bool
	sectionReset   bool
}

// CommentPrefix is an option for NewTable() that sets the comment prefix.
func CommentPrefix(v string) func(*options) error {
	return func(o *options) error { return o.setCommentPrefix(v) }
}

func (o *options) setCommentPrefix(v string) error {
	o.commentPrefix = v
	return nil
}

// EnableComments is a NewTable() option that enables comment recognition.
func EnableComments(v bool) func(*options) error {
	return func(o *options) error { return o.setEnableComments(v) }
}

func (o *options) setEnableComments(v bool) error {
	o.enableComments = v
	return nil
}

// SectionReset is a NewTable() option that enables per-section column count
// resetting. Sections are delineated by empty lines.
func SectionReset(v bool) func(*options) error {
	return func(o *options) error { return o.setSectionReset(v) }
}

func (o *options) setSectionReset(v bool) error {
	o.sectionReset = v
	return nil
}
