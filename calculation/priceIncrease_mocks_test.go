package calculation

import (
	"errors"
	"testify/mocks"
	"testify/model"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	priceIncrease     PriceIncrease
	priceProviderMock *mocks.PriceProvider
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupTest() {
	priceProviderMock := mocks.PriceProvider{}

	uts.priceIncrease = NewPriceIncrease(&priceProviderMock)
	uts.priceProviderMock = &priceProviderMock
}

func (uts *UnitTestSuite) TestCalculate() {
	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{
		{
			Timestamp: time.Now(),
			Price:     2.0,
		},
		{
			Timestamp: time.Now().Add(time.Duration(-1 * time.Minute)),
			Price:     1.0,
		},
	}, nil)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(100.0, actual)
	uts.Nil(err)
}

func (uts *UnitTestSuite) TestCalculate_Error() {
	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, nil)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.EqualError(err, "not enough data")
}

func (uts *UnitTestSuite) TestCalculate_ErrorFromPriceProvider() {
	expectedError := errors.New("oh my deuss")

	uts.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, expectedError)

	actual, err := uts.priceIncrease.PriceIncrease()

	uts.Equal(0.0, actual)
	uts.Equal(expectedError, err)
}
