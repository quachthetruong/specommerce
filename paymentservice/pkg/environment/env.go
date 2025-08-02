package environment

type Environment string

const (
	Local       Environment = "local"
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) IsProduction() bool {
	return e == Production
}
