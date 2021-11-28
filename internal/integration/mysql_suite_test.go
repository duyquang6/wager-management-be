// +build integration

package integration

import (
	"context"
	"github.com/duyquang6/wager-management-be/internal/database"
	"github.com/duyquang6/wager-management-be/internal/serverenv"
	"github.com/duyquang6/wager-management-be/internal/setup"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MySqlRepositoryTestSuite struct {
	env    *serverenv.ServerEnv
	config database.Config
	suite.Suite
}

func (p *MySqlRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	var config database.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		panic(err)
	}
	p.env = env
}

func TestMySqlRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &MySqlRepositoryTestSuite{})
}

func (p *MySqlRepositoryTestSuite) SetupTest() {
	ctx := context.Background()
	if err := p.env.Database().Migrate(ctx); err != nil {
		panic(err)
	}
}

func (p *MySqlRepositoryTestSuite) TearDownTest() {
	ctx := context.Background()
	if err := p.env.Database().MigrateDown(ctx); err != nil {
		panic(err)
	}
}
