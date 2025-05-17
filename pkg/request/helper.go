package req

import (
	"errors"
	"net/url"
	"strconv"
)

type Params struct {
	Limit  uint
	Offset uint
	Order  string
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
