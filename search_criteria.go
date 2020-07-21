package magento2

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// SearchCriteria is the query parameters that are used in th search endpoints.
// Based on:
// https://devdocs.magento.com/guides/v2.4/rest/performing-searches.html
// https://devdocs.magento.com/guides/v2.4/rest/search-endpoint.html
// https://devdocs.magento.com/guides/v2.4/rest/retrieve-filtered-responses.html
type SearchCriteria struct {
	filterGroups []SearchCriteriaFilterGroup
	sortOrders   []SearchCriteriaSortOrder
	currentPage  *int
	pageSize     *int
}

type SearchCriteriaFilterGroup struct {
	filters []SearchCriteriaFilter
}

type SearchCriteriaFilter struct {
	Field, Value, ConditionType string
}

type SearchCriteriaSortOrder struct {
	field, direction string
}

func NewSearchCriteria() *SearchCriteria {
	return &SearchCriteria{
		filterGroups: make([]SearchCriteriaFilterGroup, 0),
	}
}

func (c *SearchCriteria) AddFilterGroup(g SearchCriteriaFilterGroup) error {
	if len(g.filters) == 0 {
		return errors.New("the passed in SearchCriteriaFilterGroup doesn't have any filters set, set it using the AddFilter() method")
	}

	c.filterGroups = append(c.filterGroups, g)
	return nil
}

// SetQueryParams accepts an url.Values and add the query string params to it
func (c *SearchCriteria) SetQueryParams(v url.Values) {
	for i, g := range c.filterGroups {
		for j, f := range g.filters {
			keyFormat := "searchCriteria[filterGroups][%d][filters][%d][%s]"
			v.Set(fmt.Sprintf(keyFormat, i, j, "field"), f.Field)
			v.Set(fmt.Sprintf(keyFormat, i, j, "value"), f.Value)
			v.Set(fmt.Sprintf(keyFormat, i, j, "conditionType"), f.ConditionType)
		}
	}

	for i, s := range c.sortOrders {
		keyFormat := "searchCriteria[sortOrders][%d][%s]"
		v.Set(fmt.Sprintf(keyFormat, i, "field"), s.field)
		v.Set(fmt.Sprintf(keyFormat, i, "direction"), s.direction)
	}

	if c.currentPage != nil {
		v.Set("searchCriteria[currentPage]", strconv.Itoa(*c.currentPage))
	}
	if c.pageSize != nil {
		v.Set("searchCriteria[pageSize]", strconv.Itoa(*c.pageSize))
	}
}

func (c *SearchCriteria) AddSortOrder(s SearchCriteriaSortOrder) error {
	if s.direction != "ASC" && s.direction != "DESC" {
		return fmt.Errorf("invalid sort order direction %s", s.direction)
	}
	c.sortOrders = append(c.sortOrders, s)
	return nil
}

func (c *SearchCriteria) SetCurrentPage(i int) {
	c.currentPage = &i
}

func (c *SearchCriteria) SetPageSize(i int) {
	c.pageSize = &i
}

func NewSearchCriteriaFilterGroup() *SearchCriteriaFilterGroup {
	return &SearchCriteriaFilterGroup{
		filters: make([]SearchCriteriaFilter, 0),
	}
}

func (g *SearchCriteriaFilterGroup) AddFilter(f SearchCriteriaFilter) error {
	g.filters = append(g.filters, f)
	return nil
}
