package pkg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshrelpkg "github.com/cloudfoundry/bosh-cli/v7/release/pkg"
	. "github.com/cloudfoundry/bosh-cli/v7/release/resource"
)

func TestReg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "state/pkg")
}

func newPkg(name, fp string, deps []string) *boshrelpkg.Package {
	resource := NewResourceWithBuiltArchive(name, fp, "", "")
	return boshrelpkg.NewPackage(resource, deps)
}
