package manifest_test

import (
	"errors"

	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/v7/release/pkg/manifest"
)

var _ = Describe("NewManifestFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("parses pkg manifest successfully", func() {
		contents := `---
name: name

dependencies:
- pkg1
- pkg2

files:
- in-file1
- in-file2

excluded_files:
- ex-file1
- ex-file2
`

		err := fs.WriteFileString("/path", contents)
		Expect(err).ToNot(HaveOccurred())

		manifest, err := NewManifestFromPath("/path", fs)
		Expect(err).ToNot(HaveOccurred())
		Expect(manifest).To(Equal(Manifest{
			Name: "name",

			Dependencies: []string{"pkg1", "pkg2"},

			Files: []string{"in-file1", "in-file2"},

			ExcludedFiles: []string{"ex-file1", "ex-file2"},
		}))
	})

	It("returns error if manifest is not valid yaml", func() {
		err := fs.WriteFileString("/path", "-")
		Expect(err).ToNot(HaveOccurred())

		_, err = NewManifestFromPath("/path", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("line 1"))
	})

	It("returns error if manifest cannot be read", func() {
		err := fs.WriteFileString("/path", "-")
		Expect(err).ToNot(HaveOccurred())
		fs.ReadFileError = errors.New("fake-err")

		_, err = NewManifestFromPath("/path", fs)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("fake-err"))
	})
})
