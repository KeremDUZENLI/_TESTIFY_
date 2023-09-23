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

type MockTestSuite struct {
	suite.Suite
	priceProviderMock *mocks.PriceProvider
	priceIncrease     PriceIncrease
	myErr             error
}

func TestMockTestSuite(t *testing.T) {
	suite.Run(t, &MockTestSuite{})
}

func (mTS *MockTestSuite) SetupTest() {
	fakeDBStruct := mocks.PriceProvider{}

	mTS.priceIncrease = NewPriceIncrease(&fakeDBStruct)
	mTS.priceProviderMock = &fakeDBStruct
}

func (uTS *MockTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestCalculate_ErrorFromPriceProvider" {
		uTS.myErr = errors.New("FAIL")
	}
}

func (mTS *MockTestSuite) TestCalculate() {
	mTS.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{
		{
			Timestamp: time.Now(),
			Price:     2.0,
		},
		{
			Timestamp: time.Now().Add(time.Duration(-1 * time.Minute)),
			Price:     1.0,
		},
	}, nil)

	actual, err := mTS.priceIncrease.PriceIncrease()

	mTS.Equal(100.0, actual)
	mTS.Nil(err)
}

func (mTS *MockTestSuite) TestCalculate_Error() {
	mTS.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, nil)

	actual, err := mTS.priceIncrease.PriceIncrease()

	mTS.Equal(0.0, actual)
	mTS.EqualError(err, "not enough data")
}

func (mTS *MockTestSuite) TestCalculate_ErrorFromPriceProvider() {
	mTS.priceProviderMock.On("List", mock.Anything).Return([]*model.TimeAndPrice{}, mTS.myErr)

	actual, err := mTS.priceIncrease.PriceIncrease()

	mTS.Equal(0.0, actual)
	mTS.EqualError(err, mTS.myErr.Error())
}
