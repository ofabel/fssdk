package contract

type Orphans string

const (
	Orphans_Download Orphans = "download"
	Orphans_Delete   Orphans = "delete"
	Orphans_Ignore   Orphans = "ignore"
)

type Config struct {
	Source  string
	Target  string
	Orphans Orphans
	Include []string
	Exclude []string
	Run     []string
}
