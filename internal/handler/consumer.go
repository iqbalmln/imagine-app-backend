// Package handler
package handler

import (
	"context"

	"gitlab.com/go_graphql/internal/appctx"
	"gitlab.com/go_graphql/internal/consts"
	uContract "gitlab.com/go_graphql/internal/ucase/contract"
	"gitlab.com/go_graphql/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
