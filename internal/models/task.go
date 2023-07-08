package models

func NewTask(
	Path string,
	Name string,
	URL string,
	Dir string,
	MD5 string,
) Task {
	return Task{
		Path: Path,
		Name: Name,
		URL:  URL,
		Dir:  Dir,
		MD5:  MD5,
	}
}

type Task struct {
	Path string
	Name string
	URL  string
	Dir  string
	MD5  string
}
