package version

import (
	"gorm.io/gorm"
)

// GetLatestDBVersion retrieves the latest applied database version.
func GetLatestDBVersion(db *gorm.DB) (int, error) {
	var version GooseDBVersion
	if err := db.Order("version_id desc").Where("is_applied = ?", true).First(&version).Error; err != nil {
		return 0, err
	}
	return version.VersionID, nil
}
