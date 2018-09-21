package influx

import (
	"net/url"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
)

// Do actually executes the query to the specified InfluxDB instance hosted at the url argument.
func (q Query) Do(url string) (*influx.Response, error) {
	query := q.Build()
	response, err := executeQuery(url, query)
	if err != nil {
		return nil, errors.Wrap(err, "fail to execute query "+query)
	}
	return response, nil
}

func executeQuery(url, queryString string) (*influx.Response, error) {
	client, dbName, err := newClient(url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create InfluxDB client")
	}
	defer client.Close()

	response, err := client.Query(influx.Query{
		Command:  queryString,
		Database: dbName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to query InfluxDB")
	}
	if response.Error() != nil {
		return nil, errors.Wrap(response.Error(), "error while fetching data in InfluxDB")
	}

	return response, nil
}

type influxInfo struct {
	host             string
	user             string
	password         string
	database         string
	connectionString string
}

func newClient(url string) (c influx.Client, dbName string, err error) {
	infos, err := parseConnectionString(url)
	if err != nil {
		return nil, "", errors.Wrap(err, "invalid connection string")
	}
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:      infos.host,
		Username:  infos.user,
		Password:  infos.password,
		UserAgent: "Scalingo Utils",
	})

	if err != nil {
		return nil, "", errors.Wrap(err, "fail to create InfluxDB client")
	}

	return client, infos.database, err
}

func parseConnectionString(con string) (*influxInfo, error) {
	url, err := url.Parse(con)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse connection string")
	}

	var user, password string
	if url.User != nil {
		password, _ = url.User.Password()
		user = url.User.Username()
	}

	return &influxInfo{
		host:             url.Scheme + "://" + url.Host,
		user:             user,
		password:         password,
		database:         url.Path[1:],
		connectionString: con,
	}, nil
}
