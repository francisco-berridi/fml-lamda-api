package services

import "fmt"

func init() {

	OnBoardingService, err := NewOnBoardingService()
	fmt.Println(OnBoardingService, err)
}

func NewOnBoardingService() (*onBoardingService, error) {
	return new(onBoardingService), nil
}

func NewDataService() *DynamoDataService {
	return new(DynamoDataService)
}
