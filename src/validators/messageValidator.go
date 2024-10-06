package validators

import (
	"fmt"
	"github.com/google/uuid"
	"message-automation/src/models/base"
	"strconv"
)

func ValidateRetrieveSentMessages(limitQuery, messageId string) error {
	var err error
	limit, err := strconv.Atoi(limitQuery)
	if limitQuery != "" && err != nil {
		return &base.BadRequestError{
			Message: fmt.Sprintf("Invalid integer value for 'limit': %s", limitQuery),
		}
	} else if limit < 0 {
		return &base.BadRequestError{
			Message: fmt.Sprint("Limit value cannot be lower than 0"),
		}
	}
	if _, err = uuid.Parse(messageId); messageId != "" && err != nil {
		return &base.BadRequestError{
			Message: fmt.Sprint("messageId cannot parsed to UUID format"),
		}
	}

	return nil
}

func ValidateHandleAutomation(isActiveQuery string, currentIsActive bool) error {
	isActive, err := strconv.ParseBool(isActiveQuery)
	if err != nil {
		return &base.BadRequestError{
			Message: fmt.Sprintf("Invalid boolean value for 'isActive': %s", isActiveQuery),
		}
	}

	if currentIsActive == isActive {
		return &base.BadRequestError{
			Message: "Automation is already in the requested state",
		}
	}

	return nil
}
