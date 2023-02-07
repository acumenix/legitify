// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"context"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/analyzers/skippers"
	"github.com/Legit-Labs/legitify/internal/clients/github"
	"github.com/Legit-Labs/legitify/internal/clients/gitlab"
	"github.com/Legit-Labs/legitify/internal/collectors"
	"github.com/Legit-Labs/legitify/internal/collectors/collectors_manager"
	github2 "github.com/Legit-Labs/legitify/internal/collectors/github"
	gitlab2 "github.com/Legit-Labs/legitify/internal/collectors/gitlab"
	"github.com/Legit-Labs/legitify/internal/common/namespace"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/Legit-Labs/legitify/internal/enricher"
)

// Injectors from inject_github.go:

func setupGitHub(analyzeArgs2 *args) (*analyzeExecutor, error) {
	client, err := provideGitHubClient(analyzeArgs2)
	if err != nil {
		return nil, err
	}
	context, err := provideContext(client)
	if err != nil {
		return nil, err
	}
	v := provideGitHubCollectors(context, client, analyzeArgs2)
	collectorManager := collectors_manager.NewCollectorsManager(v)
	enginer, err := provideOpa(analyzeArgs2)
	if err != nil {
		return nil, err
	}
	skipper := skippers.NewSkipper(context)
	analyzer := analyzers.NewAnalyzer(context, enginer, skipper)
	enricherManager := enricher.NewEnricherManager(context)
	outputer := provideOutputer(context, analyzeArgs2)
	cmdAnalyzeExecutor := initializeAnalyzeExecutor(collectorManager, analyzer, enricherManager, outputer)
	return cmdAnalyzeExecutor, nil
}

// Injectors from inject_gitlab.go:

func setupGitLab(analyzeArgs2 *args) (*analyzeExecutor, error) {
	client, err := provideGitLabClient(analyzeArgs2)
	if err != nil {
		return nil, err
	}
	context, err := provideContext(client)
	if err != nil {
		return nil, err
	}
	v := provideGitLabCollectors(context, client, analyzeArgs2)
	collectorManager := collectors_manager.NewCollectorsManager(v)
	enginer, err := provideOpa(analyzeArgs2)
	if err != nil {
		return nil, err
	}
	skipper := skippers.NewSkipper(context)
	analyzer := analyzers.NewAnalyzer(context, enginer, skipper)
	enricherManager := enricher.NewEnricherManager(context)
	outputer := provideOutputer(context, analyzeArgs2)
	cmdAnalyzeExecutor := initializeAnalyzeExecutor(collectorManager, analyzer, enricherManager, outputer)
	return cmdAnalyzeExecutor, nil
}

// inject_github.go:

func provideGitHubCollectors(ctx context.Context, client *github.Client, analyzeArgs2 *args) []collectors.Collector {
	type newCollectorFunc func(ctx context.Context, client *github.Client) collectors.Collector
	var collectorsMapping = map[namespace.Namespace]newCollectorFunc{namespace.Repository: github2.NewRepositoryCollector, namespace.Organization: github2.NewOrganizationCollector, namespace.Member: github2.NewMemberCollector, namespace.Actions: github2.NewActionCollector, namespace.RunnerGroup: github2.NewRunnersCollector}

	var result []collectors.Collector
	for _, ns := range analyzeArgs2.Namespaces {
		result = append(result, collectorsMapping[ns](ctx, client))
	}

	return result
}

func provideGitHubClient(analyzeArgs2 *args) (*github.Client, error) {
	ctx := context_utils.NewContextWithSimulatedSecondaryRateLimit(context.Background(), analyzeArgs2.SimulateSecondaryRateLimit)
	return github.NewClient(ctx, analyzeArgs2.Token, analyzeArgs2.Endpoint, analyzeArgs2.
		Organizations)
}

// inject_gitlab.go:

func provideGitLabCollectors(ctx context.Context, client *gitlab.Client, analyzeArgs2 *args) []collectors.Collector {
	var collectorsMapping = map[namespace.Namespace]func(ctx context.Context, client *gitlab.Client) collectors.Collector{namespace.Organization: gitlab2.NewGroupCollector, namespace.Repository: gitlab2.NewRepositoryCollector, namespace.Member: gitlab2.NewUserCollector}

	var result []collectors.Collector
	for _, ns := range analyzeArgs2.Namespaces {
		if creator, ok := collectorsMapping[ns]; ok {
			result = append(result, creator(ctx, client))
		}
	}

	return result
}

func provideGitLabClient(analyzeArgs2 *args) (*gitlab.Client, error) {
	return gitlab.NewClient(context.Background(), analyzeArgs2.Token, analyzeArgs2.Endpoint, analyzeArgs2.Organizations)
}
