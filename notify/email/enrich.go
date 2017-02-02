package main

import (
	"fmt"
	"time"

	"remy.io/memoiz/uuid"
)

func enrichEmailing() error {
	// TODO(remy): pick maximum 2 of its notes, not recently sent to him, for which we can find content
	// TODO(remy): build an email using all these information and the enriched template
	// TODO(remy): send it this email

	var uids uuid.UUIDs
	var err error

	if uids, err = getOwners(CategoryEnrichedEmail, time.Hour*24*2, 1); err != nil {
		return err
	}

	fmt.Println(uids)

	return nil
}
