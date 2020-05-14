package aws

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/cshep4/kripto/services/data-storer/internal/model"
	"net/http"
)

type (
	Servicer interface {
		Register(ctx context.Context, user model.User) error
	}

	Handler struct {
		service Servicer
	}
)

func New(service Servicer) (*Handler, error) {
	if service == nil {
		return nil, errors.New("service_is_nil")
	}

	return &Handler{
		service: service,
	}, nil
}

func (h *Handler) Get(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func (h *Handler) Store(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
