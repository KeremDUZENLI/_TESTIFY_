package mock

import (
	"github.com/stretchr/testify/mock"
)

type Provider struct {
	mock.Mock
}

/*
func (_m *Provider) Latest() (*model.Data, error) {
	ret := _m.Called()

	var r0 *model.Data

}
*/
