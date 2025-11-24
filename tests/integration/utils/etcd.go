package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/link/v3/models"
)

func CleanupEtcdData(t *testing.T) {
	t.Helper()

	client, err := etcd.ClientFromEnv()
	require.NoError(t, err)

	dir := models.EtcdLinkDirectory

	// t.Context() has already been canceled since we're in a cleanup function.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Logf("Removing directory %s from etcd", dir)
	_, err = client.KV.Delete(ctx, dir, etcdv3.WithPrefix())
	require.NoError(t, err)
}
