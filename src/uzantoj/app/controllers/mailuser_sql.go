package controllers

import (
	"github.com/revel/revel"
)

func isMyDomain(domain_id, user_id int) bool {
	_, err := domainAssignedTo(domain_id, user_id)
	if err != nil {
		revel.AppLog.Error("Lib.IsMyDomain - Error - ", err)
	}
	return err == nil
}
