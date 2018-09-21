package influx

import (
	"fmt"
	"strings"
	"time"
)

// InfluxDB comparison operators
const (
	LessThan    = "<"
	LessOrEqual = "<="
	Equal       = "="
	MoreThan    = ">"
	MoreOrEqual = ">="
	Different   = "!="
)

// InfluxDB fill options
const (
	None     = "none"
	Null     = "null"
	Previous = "previous"
	Linear   = "linear"
)

// Query is the main structure for the query builder. All methods apply to it.
type Query struct {
	measurement    string
	conditions     condition
	fields         []field
	groupByTime    string
	groupByTag     []string
	groupByFill    string
	orderDirection string
	limit          int
}

type condition struct {
	tag        string
	comparison string
	value      string
	next       *conditionOperator
}

type field struct {
	name              string
	aggregationMethod string
}

type conditionOperator struct {
	operator  string
	condition condition
}

// NewQuery makes a new Query object. You MUST use this method to instantiate a new Query
// structure.
func NewQuery() Query {
	return Query{
		limit: -1,
	}
}

// String returns the parameter surround with single quote. You should use this method when
// dealing with values of conditions (third parameter of Where, And and Or methods).
func String(param string) string {
	return fmt.Sprintf("'%s'", param)
}

// On sets the measurement of the current query. Calling it twice will take the latest measurement
// provided.
func (q Query) On(measurement string) Query {
	q.measurement = measurement
	return q
}

// Field adds the given field to the list of fields with the given aggregation method applied. It
// is possible to add multiple fields with the same name but is highly discouraged.
func (q Query) Field(fieldname, aggregationMethod string) Query {
	q.fields = append(q.fields, field{
		name:              fieldname,
		aggregationMethod: aggregationMethod,
	})
	return q
}

// OrderByTime modifies the sorting of the result. By default, InfluxDB returns results in ascending
// time order; the first point returned has the oldest timestamp and the last point returned has the
// most recent timestamp. Calling this method with "DESC" reverses that order such that InfluxDB
// returns the points with the most recent timestamps first.
func (q Query) OrderByTime(direction string) Query {
	q.orderDirection = direction
	return q
}

// Where adds a condition on tag to the query. The comparison argument must be one of the
// constants defined in this package. Where must be called for the first condition. Calling it
// twice would remove all the previously registered conditions with And and Or.
//
// The value parameter must be surrounded with single quote if it does not represent a number. You
// can use the influx.String method to add these.
func (q Query) Where(tag, comparison, value string) Query {
	q.conditions = condition{
		tag:        tag,
		comparison: comparison,
		value:      value,
		next:       nil,
	}
	return q
}

func (q Query) addCondition(operator string, c condition) Query {
	if q.conditions == (condition{}) {
		q.conditions = c
	} else {
		lastCondition := &q.conditions
		for lastCondition.next != nil {
			lastCondition = &lastCondition.next.condition
		}
		lastCondition.next = &conditionOperator{
			operator:  operator,
			condition: c,
		}
	}
	return q
}

// And add a condition to the query separated from the previous with AND.
//
// The value parameter must be surrounded with single quote if it does not represent a number. You
// can use the influx.String method to add these.
func (q Query) And(tag, comparison, value string) Query {
	return q.addCondition("AND", condition{
		tag:        tag,
		comparison: comparison,
		value:      value,
		next:       nil,
	})
}

// Or add a condition to the query separated from the previous with OR.
//
// The value parameter must be surrounded with single quote if it does not represent a number. You
// can use the influx.String method to add these.
func (q Query) Or(tag, comparison, value string) Query {
	return q.addCondition("OR", condition{
		tag:        tag,
		comparison: comparison,
		value:      value,
		next:       nil,
	})
}

// Limit sets the limit of the current query. Calling it twice will take the latest limit
// provided.
// It limits the number of points returned by the query.
func (q Query) Limit(limit int) Query {
	q.limit = limit
	return q
}

// GroupByTime groups query results by a time interval. Calling it twice will take the latest
// duration provided.
func (q Query) GroupByTime(duration time.Duration) Query {
	q.groupByTime = duration.String()
	return q
}

// GroupByTag groups query results by a user-specified set of tags. Every call to this method adds
// the given slice of tags to the existing set of tags.
func (q Query) GroupByTag(tag ...string) Query {
	q.groupByTag = append(q.groupByTag, tag...)
	return q
}

// Fill sets the fill behaviour of the current query. Calling it twice will take the latest fill
// provided.
// It changes the value reported for time intervals that have no data.
func (q Query) Fill(value string) Query {
	q.groupByFill = value
	return q
}

// Build constructs the InfluxQL query in a string form.
func (q Query) Build() string {
	query := "SELECT"

	for i, f := range q.fields {
		if i != 0 {
			query += ","
		}
		query += fmt.Sprintf(" %s(\"%s\") AS \"%s\"", f.aggregationMethod, f.name, f.name)
	}

	query += fmt.Sprintf(" FROM %s", q.measurement)
	if q.conditions != (condition{}) {
		query += fmt.Sprintf(" WHERE %s", q.conditions.build())
	}

	if q.groupByTime != "" || len(q.groupByTag) > 0 {
		query += " GROUP BY "
		if q.groupByTime != "" {
			query += fmt.Sprintf("time(%s)", q.groupByTime)
			if len(q.groupByTag) > 0 {
				query += ","
			}
		}

		if len(q.groupByTag) > 0 {
			query += strings.Join(q.groupByTag, ",")
		}

		if q.groupByFill != "" {
			query += fmt.Sprintf(" fill(%s)", q.groupByFill)
		}
	}

	if q.orderDirection != "" {
		query += fmt.Sprintf(" ORDER BY time %s", q.orderDirection)
	}

	if q.limit != -1 {
		query += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	return query
}

func (c condition) build() string {
	query := fmt.Sprintf("%s %s %s", c.tag, c.comparison, c.value)
	if c.next != nil {
		query += fmt.Sprintf(" %s %s", c.next.operator, c.next.condition.build())
	}
	return query
}
