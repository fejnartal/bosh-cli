package director_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/v7/director"
)

var _ = Describe("NewOSVersionSlug", func() {
	It("populates slug with os and version", func() {
		slug := NewOSVersionSlug("name", "ver")
		Expect(slug.OS()).To(Equal("name"))
		Expect(slug.Version()).To(Equal("ver"))
	})

	It("panics if name is empty", func() {
		Expect(func() { NewOSVersionSlug("", "") }).To(Panic())
	})

	It("panics if version is empty", func() {
		Expect(func() { NewOSVersionSlug("name", "") }).To(Panic())
	})
})

var _ = Describe("OSVersionSlug", func() {
	Describe("String", func() {
		It("returns os/version", func() {
			Expect(NewOSVersionSlug("os", "ver").String()).To(Equal("os/ver"))
		})
	})

	Describe("IsProvided", func() {
		It("returns true if name and id are specified", func() {
			Expect(NewOSVersionSlug("os", "version").IsProvided()).To(BeTrue())
		})

		It("returns false if it's empty", func() {
			Expect(OSVersionSlug{}.IsProvided()).To(BeFalse())
		})
	})

	Describe("UnmarshalFlag", func() {
		var (
			slug *OSVersionSlug
		)

		BeforeEach(func() {
			slug = &OSVersionSlug{}
		})

		It("populates slug", func() {
			err := slug.UnmarshalFlag("os/ver")
			Expect(err).ToNot(HaveOccurred())
			Expect(*slug).To(Equal(NewOSVersionSlug("os", "ver")))
		})

		It("returns an error if string doesnt have 2 pieces", func() {
			err := slug.UnmarshalFlag("1")
			Expect(err).To(Equal(errors.New("Expected OS '1' to be in format 'name/version'")))

			err = slug.UnmarshalFlag("1.2.3")
			Expect(err).To(Equal(errors.New("Expected OS '1.2.3' to be in format 'name/version'")))
		})

		It("returns an error if name is empty", func() {
			err := slug.UnmarshalFlag("/ver")
			Expect(err).To(Equal(errors.New("Expected OS '/ver' to specify non-empty name")))
		})

		It("returns an error if version is empty", func() {
			err := slug.UnmarshalFlag("name/")
			Expect(err).To(Equal(errors.New("Expected OS 'name/' to specify non-empty version")))
		})
	})

	Describe("UnmarshalJSON", func() {
		var (
			slug *OSVersionSlug
		)

		BeforeEach(func() {
			slug = &OSVersionSlug{}
		})

		It("populates slug", func() {
			err := slug.UnmarshalJSON([]byte(`"os/ver"`))
			Expect(err).ToNot(HaveOccurred())
			Expect(*slug).To(Equal(NewOSVersionSlug("os", "ver")))
		})

		It("returns an error if string doesnt have 2 pieces", func() {
			err := slug.UnmarshalJSON([]byte(`"os"`))
			Expect(err).To(Equal(errors.New("Expected OS 'os' to be in format 'name/version'")))

			err = slug.UnmarshalJSON([]byte(`"os/2/3"`))
			Expect(err).To(Equal(errors.New("Expected OS 'os/2/3' to be in format 'name/version'")))
		})

		It("returns an error if os is empty", func() {
			err := slug.UnmarshalJSON([]byte(`"/ver"`))
			Expect(err).To(Equal(errors.New("Expected OS '/ver' to specify non-empty name")))
		})

		It("returns an error if version is empty", func() {
			err := slug.UnmarshalJSON([]byte(`"os/"`))
			Expect(err).To(Equal(errors.New("Expected OS 'os/' to specify non-empty version")))
		})
	})
})
