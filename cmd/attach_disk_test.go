package cmd_test

import (
	"errors"

	. "github.com/cloudfoundry/bosh-cli/v7/cmd"
	. "github.com/cloudfoundry/bosh-cli/v7/cmd/opts"
	boshdir "github.com/cloudfoundry/bosh-cli/v7/director"
	fakedir "github.com/cloudfoundry/bosh-cli/v7/director/directorfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AttachDisk", func() {
	var (
		command    AttachDiskCmd
		deployment *fakedir.FakeDeployment
	)

	BeforeEach(func() {
		deployment = &fakedir.FakeDeployment{}

		command = NewAttachDiskCmd(deployment)
	})

	Describe("Run", func() {
		var (
			opts           AttachDiskOpts
			act            func() error
			instanceSlug   boshdir.InstanceSlug
			diskCid        string
			diskProperties string
		)

		BeforeEach(func() {
			act = func() error {
				err := command.Run(opts)
				return err
			}

			instanceSlug = boshdir.NewInstanceSlug("instance-group-name", "1")
			diskCid = "some-disk-id"
			diskProperties = "copy"

			opts = AttachDiskOpts{
				Args: AttachDiskArgs{
					Slug:    instanceSlug,
					DiskCID: diskCid,
				},
			}
			opts.DiskProperties = diskProperties
		})

		It("tells the director to attach a disk", func() {
			err := act()
			Expect(err).NotTo(HaveOccurred())
			Expect(deployment.AttachDiskCallCount()).To(Equal(1))

			receivedInstanceSlug, receivedDiskCid, receivedDiskProperties := deployment.AttachDiskArgsForCall(0)

			Expect(receivedInstanceSlug).To(Equal(instanceSlug))
			Expect(receivedDiskCid).To(Equal(diskCid))
			Expect(receivedDiskProperties).To(Equal("copy"))
		})

		Context("attaching a disk returns an error", func() {

			BeforeEach(func() {
				deployment.AttachDiskReturns(errors.New("director returned an error attaching a disk"))
			})

			It("Should return an error if director attaching a disk fails", func() {
				err := act()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("director returned an error attaching a disk"))
			})
		})
	})
})
