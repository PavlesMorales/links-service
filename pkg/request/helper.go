package req

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

type Params struct {
	Limit  uint
	Offset uint
	Order  string
}

const (
	GroupByMonth = "month"
	GroupByDay   = "day"
)

type StatParams struct {
	From time.Time
	To   time.Time
	By   string
}

func newBy(by string) (string, error) {
	if by != GroupByMonth && by != GroupByDay {
		return "", errors.New("invalid param by must be a day or month")
	}
	return by, nil
}

func timeParse(val string) (time.Time, error) {
	return time.Parse("2006-01-02", val)
}

func GetStatParams(queryParams url.Values) (*StatParams, error) {

	by, err := newBy(queryParams.Get("by"))
	if err != nil {
		return nil, err
	}
	to, err := timeParse(queryParams.Get("to"))
	if err != nil {
		return nil, err
	}
	from, err := timeParse(queryParams.Get("from"))
	if err != nil {
		return nil, err
	}

	return &StatParams{
		From: from,
		To:   to,
		By:   by,
	}, nil
}

func GetParams(queryParams url.Values) (*Params, error) {
	limit, err := strconv.ParseUint(queryParams.Get("limit"), 10, 32)
	if err != nil {
		return nil, errors.New("limit parameter must be an integer")
	}
	offset, err := strconv.ParseUint(queryParams.Get("offset"), 10, 32)
	if err != nil {
		return nil, errors.New("offset parameter must be an integer")
	}

	order := queryParams.Get("order")

	if order == "" || order != "asc" && order != "desc" {
		order = "asc"
	}

	return &Params{
		Limit:  uint(limit),
		Offset: uint(offset),
		Order:  order,
	}, nil

}
