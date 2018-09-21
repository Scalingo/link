package influx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConditionBuilder(t *testing.T) {
	examples := []struct {
		Name      string
		condition condition
		Expected  string
	}{
		{
			Name: "a simple condition",
			condition: condition{
				tag:        "tag",
				comparison: Equal,
				value:      `"test"`,
			},
			Expected: `tag = "test"`,
		}, {
			Name: "multiple conditions",
			condition: condition{
				tag:        "tag",
				comparison: Different,
				value:      `"hi"`,
				next: &conditionOperator{
					operator: "AND",
					condition: condition{
						tag:        "time",
						comparison: LessThan,
						value:      "now() - 3m",
					},
				},
			},
			Expected: `tag != "hi" AND time < now() - 3m`,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			assert.Equal(t, example.Expected, example.condition.build())
		})
	}
}

func TestQueryBuilder(t *testing.T) {
	examples := []struct {
		Name     string
		Query    Query
		Expected string
	}{
		{
			Name:     "a basic query",
			Query:    NewQuery().On("serie").Field("test", "mean"),
			Expected: `SELECT mean("test") AS "test" FROM serie`,
		}, {
			Name:     "with a basic condition",
			Query:    NewQuery().On("serie").Field("test", "mean").Where("condF", Equal, `"value"`),
			Expected: `SELECT mean("test") AS "test" FROM serie WHERE condF = "value"`,
		}, {
			Name:     "with a complex condition",
			Query:    NewQuery().On("serie").Field("test", "mean").Where("cond1F", Equal, `"value"`).And("cond2F", MoreOrEqual, "time() - 3m"),
			Expected: `SELECT mean("test") AS "test" FROM serie WHERE cond1F = "value" AND cond2F >= time() - 3m`,
		}, {
			Name:     "with an even more complex condition",
			Query:    NewQuery().On("serie").Field("test", "mean").Where("cond1F", Equal, `"value"`).And("cond2F", MoreOrEqual, "time() - 3m").And("cond3F", Equal, "'app_id'"),
			Expected: `SELECT mean("test") AS "test" FROM serie WHERE cond1F = "value" AND cond2F >= time() - 3m AND cond3F = 'app_id'`,
		}, {
			Name:     "with a group by tag",
			Query:    NewQuery().On("serie").Field("f1", "last").GroupByTag("tag1"),
			Expected: `SELECT last("f1") AS "f1" FROM serie GROUP BY tag1`,
		}, {
			Name:     "with multiple group by time",
			Query:    NewQuery().On("serie").Field("f1", "last").GroupByTag("tag1").GroupByTag("tag2"),
			Expected: `SELECT last("f1") AS "f1" FROM serie GROUP BY tag1,tag2`,
		}, {
			Name:     "with a group by time",
			Query:    NewQuery().On("serie").Field("f1", "last").GroupByTime(10 * time.Minute),
			Expected: `SELECT last("f1") AS "f1" FROM serie GROUP BY time(10m0s)`,
		}, {
			Name:     "with time, tags and fill",
			Query:    NewQuery().On("serie").Field("f1", "last").GroupByTag("tag1").GroupByTag("tag2").GroupByTime(1 * time.Second).Fill(Previous),
			Expected: `SELECT last("f1") AS "f1" FROM serie GROUP BY time(1s),tag1,tag2 fill(previous)`,
		}, {
			Name:     "with an order by",
			Query:    NewQuery().On("serie").Field("f1", "mean").OrderByTime("DESC"),
			Expected: `SELECT mean("f1") AS "f1" FROM serie ORDER BY time DESC`,
		}, {
			Name:     "with a limit",
			Query:    NewQuery().On("serie").Field("f1", "mean").Limit(1),
			Expected: `SELECT mean("f1") AS "f1" FROM serie LIMIT 1`,
		},
	}
	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			assert.Equal(t, example.Expected, example.Query.Build())
		})
	}

	t.Run("query builder should be stateless", func(t *testing.T) {
		query := NewQuery()
		// Calling all the methods should not modify the receiving object.
		query.On("serie").Field("f1", "mean").
			Where("cond1F", Equal, `"value"`).And("cond2F", MoreOrEqual, "time() - 3m").
			GroupByTag("tag1").GroupByTime(1 * time.Second).Fill(Previous).
			OrderByTime("DESC").Limit(1)
		assert.Equal(t, NewQuery().Build(), query.Build())
	})
}

func TestString(t *testing.T) {
	t.Run("it should add surrounding single quotes around the parameter", func(t *testing.T) {
		assert.Equal(t, "'biniou'", String("biniou"))
	})
}
