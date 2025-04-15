package utils

import (
	"log"
	"strconv"
	"strings"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
	"github.com/gin-gonic/gin"
)

func FilterFromContext(
	c *gin.Context,
) (*dtos.Filter, error) {
	filter := &dtos.Filter{}

	page, per, err := paginationFromContext(c)
	if err != nil {
		return filter, err
	}

	filter.Page = page
	filter.Per = per
	filter.From = strings.TrimSpace(c.Query("from"))
	filter.To = strings.TrimSpace(c.Query("to"))
	filter.Status = strings.TrimSpace(c.Query("status"))
	filter.Type = strings.TrimSpace(c.Query("type"))
	filter.Token = strings.TrimSpace(c.Query("token"))
	filter.ResetKey = strings.TrimSpace(c.Query("reset_key"))
	filter.Term = strings.TrimSpace(c.Query("term"))
	filter.UUID = strings.TrimSpace(c.Query("uuid"))
	filter.Provider = null.NullValue(strings.TrimSpace(c.Query("provider")))
	filter.WalletType = strings.TrimSpace(c.Query("wallet_type"))
	filter.TransactionCategory = strings.TrimSpace(c.Query("transaction_category"))
	filter.BusinessAccountID = null.NullValue[string](strings.TrimSpace(c.Query("business_account_id")))
	filter.OrganizationID = null.NullValue[string](strings.TrimSpace(c.Query("organization_id")))

	isValid := strings.TrimSpace(c.Query("valid"))
	if isValid != "" {
		isValid, err := strconv.ParseBool(isValid)
		if err != nil {
			log.Printf("Failed to parse valid query param [%v]", isValid)
			return filter, apperrs.New(
				err,
				apperrs.BadRequest,
			)
		}

		filter.SetValid(isValid)
	}

	return filter, nil
}

func paginationFromContext(
	c *gin.Context,
) (int, int, error) {
	page := 1
	per := 20

	var err error

	pageQuery := strings.TrimSpace(c.Query("page"))
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			log.Printf("Failed to parse page query param [%v]", pageQuery)
			return page, per, apperrs.New(
				err,
				apperrs.BadRequest,
			)
		}
	}

	perQuery := strings.TrimSpace(c.Query("per"))
	if perQuery != "" {
		per, err = strconv.Atoi(perQuery)
		if err != nil {
			log.Printf("Failed to parse per query param [%v]", perQuery)
			return page, per, apperrs.New(
				err,
				apperrs.BadRequest,
			)
		}
	}

	return page, per, nil
}
