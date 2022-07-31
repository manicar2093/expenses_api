package repos_test

import (
	"context"
	"testing"

	"github.com/manicar2093/expenses_api/internal/connections"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	conn = connections.GetMongoConn()
)

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
	conn.Drop(context.Background()) //nolint: errcheck
}
