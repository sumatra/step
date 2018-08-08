package mocks

import "github.com/coinbase/step/aws"

type MockClients struct {
	S3     *MockS3Client
	Lambda *MockLambdaClient
	SFN    *MockSFNClient
}

func (awsc *MockClients) S3Client(*string, *string, *string) aws.S3API {
	return awsc.S3
}

func (awsc *MockClients) LambdaClient(*string, *string, *string) aws.LambdaAPI {
	return awsc.Lambda
}

func (awsc *MockClients) SFNClient(*string, *string, *string) aws.SFNAPI {
	return awsc.SFN
}

func MockAwsClients() *MockClients {
	return &MockClients{
		&MockS3Client{},
		&MockLambdaClient{},
		&MockSFNClient{},
	}
}
