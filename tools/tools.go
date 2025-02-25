//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "golang.org/x/tools/cmd/goimports"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
