package scraper

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	// "github.com/qor/l10n"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Topic struct {
	gorm.Model
	// l10n.Locale
	sorting.Sorting
	Name string
	Code string

	Topics  []Topic
	TopicID uint
}

func (topic Topic) Validate(db *gorm.DB) {
	if strings.TrimSpace(topic.Name) == "" {
		db.AddError(validations.NewError(topic, "Name", "Name can not be empty"))
	}
}

func (topic Topic) DefaultPath() string {
	if len(topic.Code) > 0 {
		return fmt.Sprintf("/topic/%s", topic.Code)
	}
	return "/"
}
