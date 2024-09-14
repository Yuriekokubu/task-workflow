package version

// GooseDBVersion represents the database schema for goose_db_version.
type GooseDBVersion struct {
	ID        int
	VersionID int
	IsApplied bool
	Tstamp    string
}

// TableName returns the table name for GooseDBVersion.
func (GooseDBVersion) TableName() string {
	return "goose_db_version"
}
