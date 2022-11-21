package repos_test

import (
	"testing"

	"github.com/manicar2093/expenses_api/internal/connections"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	conn = connections.GetGormConnection()
)

func TestRepos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repos Suite")
}
