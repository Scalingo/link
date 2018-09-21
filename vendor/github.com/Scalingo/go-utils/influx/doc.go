// Package influx provides a query builder to ease the writing of query to InfluxDB with the
// InfluxQL language.
//
// After instantiating a new query object with NewQuery, call the different methods to adapt the
// query to fit your needs. You can call Build to get the InfluxQL query in a string form.
// Eventually call Do to actually execute the query.
//
// Example Usage
//
//   package models
//
//   func (a *Alert) InfluxQuery() influx.Query {
//    	query := influx.NewQuery()
//
//    	query = query.Where("time", influx.LessThan, "now() - 1m").And("time", influx.MoreThan, "now() - 5m")
//    	query = query.And("application_id", influx.Equal, influx.String(a.AppID))
//    	query = query.GroupByTime(5*time.Minute).Fill(influx.None).OrderBy("time", "DESC").Limit(1)
//
//    	if a.IsRouterMetric() {
//    		query = query.On("router").Field("value", "mean")
//    		if a.Metric == metric.MetricRouterP95ResponseTime {
//    			query = query.And("type", influx.Equal, influx.String("p95"))
//    		}
//    	} else {
//    		query.And("container_type", influx.Equal, influx.String(a.ContainerType))
//    		query = query.On(a.Metric).Field("value", "mean").Field("limit", "last")
//    		query = query.GroupByTag("container_type").GroupByTag("container_index")
//    	}
//    	return query
//    }
package influx
