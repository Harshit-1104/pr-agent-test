package cart_context

import (
	"context"
	"fmt"

	dynamodbv2 "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pr-agent-test/internal/registry"
)

type LandingSourceDetails struct {
	landingSource registry.CartLandingSource
	entryPoints   []*EntryPoint
}

func EmptyLandingSourceDetails() *LandingSourceDetails {
	return &LandingSourceDetails{
		landingSource: registry.CartLandingSourceUnspecified,
	}
}

func NewLandingSourceDetails(landingSource registry.CartLandingSource, entryPoints []*EntryPoint) *LandingSourceDetails {
	return &LandingSourceDetails{landingSource: landingSource, entryPoints: entryPoints}
}

func (s *LandingSourceDetails) IsEmpty() bool {
	return s == nil || (s.landingSource == registry.CartLandingSourceUnspecified)
}

func (s *LandingSourceDetails) LandingSource() registry.CartLandingSource {
	if s == nil {
		return registry.CartLandingSourceUnspecified
	}
	return s.landingSource
}

func (s *LandingSourceDetails) EntryPoints() []*EntryPoint {
	return s.entryPoints
}

func (s *LandingSourceDetails) Clone() *LandingSourceDetails {
	if s == nil {
		return nil
	}

	var cl []*EntryPoint
	for _, e := range s.entryPoints {
		cl = append(cl, NewEntryPoint(e.EntryPointType()))
	}

	return NewLandingSourceDetails(
		s.landingSource,
		cl,
	)
}

type EntryPoint struct {
	entryPointType registry.EntryPointType
	Client         *dynamodbv2.Client
}

func EmptyEntryPoint() *EntryPoint {
	return &EntryPoint{
		entryPointType: registry.EntryPointTypeUnspecified,
	}
}

func NewEntryPoint(entryPointType registry.EntryPointType) *EntryPoint {
	return &EntryPoint{
		entryPointType: entryPointType,
	}
}

func (e *EntryPoint) EntryPointType() registry.EntryPointType {
	if e == nil {
		return registry.EntryPointTypeUnspecified
	}
	return e.entryPointType
}

func (e *EntryPoint) IsEmpty() bool {
	return e == nil || e.entryPointType == registry.EntryPointTypeUnspecified
}

func (e *EntryPoint) saveEntryPointInCache(ctx context.Context) error {
	_, err := e.Client.PutItem(
		context.Background(),
		&dynamodbv2.PutItemInput{},
	)
	if err != nil {
		fmt.Printf("%v", err)
	}

	return nil
}
