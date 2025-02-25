package manifest_test

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	biproperty "github.com/cloudfoundry/bosh-utils/property"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	fakeuuid "github.com/cloudfoundry/bosh-utils/uuid/fakes"
	"github.com/cppforlife/go-patch/patch"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshtpl "github.com/cloudfoundry/bosh-cli/v7/director/template"
	"github.com/cloudfoundry/bosh-cli/v7/installation/manifest"
	"github.com/cloudfoundry/bosh-cli/v7/installation/manifest/fakes"
	birelsetmanifest "github.com/cloudfoundry/bosh-cli/v7/release/set/manifest"
)

type manifestFixtures struct {
	validManifest             string
	missingPrivateKeyManifest string
}

var _ = Describe("Parser", func() {
	comboManifestPath := "/path/to/fake-deployment-manifest"
	releaseSetManifest := birelsetmanifest.Manifest{}
	var (
		fakeFs            *fakesys.FakeFileSystem
		fakeUUIDGenerator *fakeuuid.FakeGenerator
		parser            manifest.Parser
		logger            boshlog.Logger
		fakeValidator     *fakes.FakeValidator
		fixtures          manifestFixtures
	)
	BeforeEach(func() {
		fakeValidator = fakes.NewFakeValidator()
		fakeValidator.SetValidateBehavior([]fakes.ValidateOutput{
			{Err: nil},
		})
		fakeFs = fakesys.NewFakeFileSystem()
		logger = boshlog.NewLogger(boshlog.LevelNone)
		fakeUUIDGenerator = fakeuuid.NewFakeGenerator()
		parser = manifest.NewParser(fakeFs, fakeUUIDGenerator, logger, fakeValidator)
		fixtures = manifestFixtures{
			validManifest: `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  properties:
    fake-property-name:
      nested-property: fake-property-value
`,
			missingPrivateKeyManifest: `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    password: fake-password
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`,
		}
	})

	Describe("#Parse", func() {
		Context("when combo manifest path does not exist", func() {
			It("returns an error", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when parser fails to read the combo manifest file", func() {
			JustBeforeEach(func() {
				err := fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)
				Expect(err).ToNot(HaveOccurred())
				fakeFs.ReadFileError = errors.New("fake-read-file-error")
			})

			It("returns an error", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("with a valid manifest", func() {
			BeforeEach(func() {
				err := fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)
				Expect(err).ToNot(HaveOccurred())
			})

			It("parses installation from combo manifest", func() {
				installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).ToNot(HaveOccurred())

				Expect(installationManifest).To(Equal(manifest.Manifest{
					Name: "fake-deployment-name",
					Template: manifest.ReleaseJobRef{
						Name:    "fake-cpi-job-name",
						Release: "fake-cpi-release-name",
					},
					Properties: biproperty.Map{
						"fake-property-name": biproperty.Map{
							"nested-property": "fake-property-value",
						},
					},
					Mbus: "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
				}))
			})
		})

		Context("when ssh tunnel config is present", func() {
			Context("with raw private key", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIByQIBAAJhANs/tl5Tv7CD0Gz5TYocWZbGwHIjDU8dY1oszVMb8bhybfF4y88k
      7oaFYlyZ0oZATpx1EGXZAcgDszq5XSXhYKWQL6+u0qEylWsbra7qQefm2+WbZDfh
      ugqbt+kD0F6CjQIDAQABAmBS8yDxQShGBSjnAc9XUHCIvftzc1WGuCytokOwjOMA
      ELMN59DcNzHTTUWwmTXwOwWPnz1c7PYRnFmy99dEcyWeugU0C5QS96XWwGdXcOjY
      Kr1q/yDJZh416/nWkyGlIOECMQDvT36aXqf0xZHb47aEWmeezGS9IyK1BDMEqvcD
      DNU/GK86ymoEqtIyQbnuBUqSbkUCMQDqigydhP7j1IGABdVrWXX/WFhABjAmNWrf
      YYEecgjhjdM83QSkpwu7tYCHtZjny6kCMCZO6GpXurUxJ0823ZHEUxAVkg7A4B5w
      BKa7o30GgeBu2CYmHuCOY8WNxfC3Qh+8rQIwGQIXTkR8GTbzh/8XPpcPaea1oj4G
      rExN1PvElMZ8A/DncTnv4M6fBajYx5+pai3hAjBui9LTgI1fZeOtgBEo+Q3ZLm/O
      bX621YeY03FF5+TCF6Zwk4yT/NWMwJz8Fpb9QQA=
      -----END RSA PRIVATE KEY-----
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the raw private key field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest).To(Equal(manifest.Manifest{
							Name: "fake-deployment-name",
							Template: manifest.ReleaseJobRef{
								Name:    "fake-cpi-job-name",
								Release: "fake-cpi-release-name",
							},
							Properties: biproperty.Map{},
							Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
						}))

					})
				})
				Context("that is invalid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      no valid private key
      -----END RSA PRIVATE KEY-----
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Invalid private key for ssh tunnel"))
					})
				})
			})

			Context("with new format raw private key", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAdwAAAAdz
      c2gtcnNhAAAAAwEAAQAAAGEA2z+2XlO/sIPQbPlNihxZlsbAciMNTx1jWizNUxvx
      uHJt8XjLzyTuhoViXJnShkBOnHUQZdkByAOzOrldJeFgpZAvr67SoTKVaxutrupB
      5+bb5ZtkN+G6Cpu36QPQXoKNAAABoOJ0x/nidMf5AAAAB3NzaC1yc2EAAABhANs/
      tl5Tv7CD0Gz5TYocWZbGwHIjDU8dY1oszVMb8bhybfF4y88k7oaFYlyZ0oZATpx1
      EGXZAcgDszq5XSXhYKWQL6+u0qEylWsbra7qQefm2+WbZDfhugqbt+kD0F6CjQAA
      AAMBAAEAAABgUvMg8UEoRgUo5wHPV1BwiL37c3NVhrgsraJDsIzjABCzDefQ3Dcx
      001FsJk18DsFj589XOz2EZxZsvfXRHMlnroFNAuUEvel1sBnV3Do2Cq9av8gyWYe
      Nev51pMhpSDhAAAAMG6L0tOAjV9l462AESj5Ddkub85tfrbVh5jTcUXn5MIXpnCT
      jJP81YzAnPwWlv1BAAAAADEA709+ml6n9MWR2+O2hFpnnsxkvSMitQQzBKr3AwzV
      PxivOspqBKrSMkG57gVKkm5FAAAAMQDqigydhP7j1IGABdVrWXX/WFhABjAmNWrf
      YYEecgjhjdM83QSkpwu7tYCHtZjny6kAAAAUaW1wb3J0ZWQtb3BlbnNzaC1rZXkB
      AgMEBQYH
      -----END OPENSSH PRIVATE KEY-----
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the raw private key field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest).To(Equal(manifest.Manifest{
							Name: "fake-deployment-name",
							Template: manifest.ReleaseJobRef{
								Name:    "fake-cpi-job-name",
								Release: "fake-cpi-release-name",
							},
							Properties: biproperty.Map{},
							Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
						}))

					})
				})
				Context("that is invalid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: |
      -----BEGIN OPENSSH PRIVATE KEY-----
      no valid private key
      -----END OPENSSH PRIVATE KEY-----
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Invalid private key for ssh tunnel"))
					})
				})
			})

			Context("with private key format", func() {
				Context("that is unsupported", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: |
      ---- BEGIN SSH2 ENCRYPTED PRIVATE KEY ----
      Comment: "imported-openssh-key"
      P2/56wAAAb4AAAA3aWYtbW9kbntzaWdue3JzYS1wa2NzMS1zaGExfSxlbmNyeXB0e3JzYS
      1wa2NzMXYyLW9hZXB9fQAAAARub25lAAABbwAAAWsAAAARAQABAAAC/1LzIPFBKEYFKOcB
      z1dQcIi9+3NzVYa4LK2iQ7CM4wAQsw3n0Nw3MdNNRbCZNfA7BY+fPVzs9hGcWbL310RzJZ
      66BTQLlBL3pdbAZ1dw6NgqvWr/IMlmHjXr+daTIaUg4QAAAwDbP7ZeU7+wg9Bs+U2KHFmW
      xsByIw1PHWNaLM1TG/G4cm3xeMvPJO6GhWJcmdKGQE6cdRBl2QHIA7M6uV0l4WClkC+vrt
      KhMpVrG62u6kHn5tvlm2Q34boKm7fpA9Bego0AAAF/bovS04CNX2XjrYARKPkN2S5vzm1+
      ttWHmNNxRefkwhemcJOMk/zVjMCc/BaW/UEAAAABgOqKDJ2E/uPUgYAF1WtZdf9YWEAGMC
      Y1at9hgR5yCOGN0zzdBKSnC7u1gIe1mOfLqQAAAYDvT36aXqf0xZHb47aEWmeezGS9IyK1
      BDMEqvcDDNU/GK86ymoEqtIyQbnuBUqSbkU=
      ---- END SSH2 ENCRYPTED PRIVATE KEY ----
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Unsupported private key format for ssh tunnel"))
					})
				})
			})

			Context("with private key path", func() {
				Context("with absolute private_key path", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: /path/to/fake-ssh-key.pem
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
						err = fakeFs.WriteFileString("/path/to/fake-ssh-key.pem", "--- BEGIN KEY --- blah --- END KEY ---")
						Expect(err).ToNot(HaveOccurred())
					})

				})

				Context("with relative private_key path", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: tmp/fake-ssh-key.pem
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
						err = fakeFs.WriteFileString("/path/to/tmp/fake-ssh-key.pem", "--- BEGIN KEY --- blah --- END KEY ---")
						Expect(err).ToNot(HaveOccurred())
					})
				})

				Context("with private_key path beginning with '~'", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: ~/tmp/fake-ssh-key.pem
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
						fakeFs.ExpandPathExpanded = "/Users/foo/tmp/fake-ssh-key.pem"
						err = fakeFs.WriteFileString(fakeFs.ExpandPathExpanded, "--- BEGIN KEY --- blah --- END KEY ---")
						Expect(err).ToNot(HaveOccurred())
					})

				})

				Context("when expanding to the home directory fails", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: ~/tmp/fake-ssh-key.pem
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeFs.ExpandPathErr = errors.New("fake-expand-error")
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Expanding private_key path: fake-expand-error"))
					})
				})

				Context("when file does not exist", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: /bar/fake-ssh-key.pem
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})
					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err.Error()).To(ContainSubstring("Reading private key from /bar/fake-ssh-key.pem"))
					})
				})
			})

		})

		It("handles installation manifest validation errors", func() {
			err := fakeFs.WriteFileString(comboManifestPath, fixtures.validManifest)
			Expect(err).ToNot(HaveOccurred())

			fakeValidator.SetValidateBehavior([]fakes.ValidateOutput{
				{Err: errors.New("nope")},
			})

			_, err = parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Validating installation manifest: nope"))
		})

		Context("when interpolating variables", func() {
			BeforeEach(func() {
				fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
				fakeFs.ExpandPathExpanded = "/Users/foo/tmp/fake-ssh-key.pem"

				err := fakeFs.WriteFileString(comboManifestPath, `---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  ssh_tunnel:
    host: 54.34.56.8
    port: 22
    user: fake-ssh-user
    private_key: ((url))
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
`)
				Expect(err).ToNot(HaveOccurred())
				err = fakeFs.WriteFileString("/Users/foo/tmp/fake-ssh-key.pem", "--- BEGIN KEY --- blah --- END KEY ---")
				Expect(err).ToNot(HaveOccurred())
			})

			It("resolves their values", func() {
				vars := boshtpl.StaticVariables{"url": "~/tmp/fake-ssh-key.pem"}
				ops := patch.Ops{
					patch.ReplaceOp{Path: patch.MustNewPointerFromString("/name"), Value: "replaced-name"},
				}

				installationManifest, err := parser.Parse(comboManifestPath, vars, ops, releaseSetManifest)
				Expect(err).ToNot(HaveOccurred())

				Expect(installationManifest).To(Equal(manifest.Manifest{
					Name: "replaced-name",
					Template: manifest.ReleaseJobRef{
						Name:    "fake-cpi-job-name",
						Release: "fake-cpi-release-name",
					},
					Properties: biproperty.Map{},
					Mbus:       "http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868",
				}))
			})

			It("returns an error if variable key is missing", func() {
				_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Expected to find variables: url"))
			})
		})

		Context("when CA cert is present", func() {
			Context("with raw certificate", func() {
				Context("that is valid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  cert:
    ca: |
      -----BEGIN CERTIFICATE-----
      MIIC+TCCAeGgAwIBAgIQLzf5Fs3v+Dblm+CKQFxiKTANBgkqhkiG9w0BAQsFADAm
      MQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkwHhcNMTcwNTE2
      MTUzNTI4WhcNMTgwNTE2MTUzNTI4WjAmMQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoT
      DUNsb3VkIEZvdW5kcnkwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC+
      4E0QJMOpQwbHACvrZ4FleP4/DMFvYUBySfKzDOgd99Nm8LdXuJcI1SYHJ3sV+mh0
      +cQmRt8U2A/lw7bNU6JdM0fWHa/2nGjSBKWgPzba68NdsmwjqUjLatKpr1yvd384
      PJJKC7NrxwvChgB8ui84T4SrXHCioYMDEDIqLGmHJHMKnzQ17nu7ECO4e6QuCfnH
      RDs7dTjomTAiFuF4fh4SPgEDMGaCE5HZr4t3gvc9n4UftpcCpi+Jh+neRiWx+v37
      ZAYf2kp3wWtYDlgWk06cZzHZZ9uYZFwHDNHdDKHxGGvAh2Rm6rpPF2oA6OEyx6BH
      85/STCgSMCnV1Wkd+1yPAgMBAAGjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
      Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBGvGggx3IM4KCMpVDSv9zFKX4K
      IuCRQ6VFab3sgnlelMFaMj3+8baJ/YMko8PP1wVfUviVgKuiZO8tqL00Yo4s1WKp
      x3MLIG4eBX9pj0ZVRa3kpcF2Wvg6WhrzUzONf7pfuz/9avl77o4aSt4TwyCvM4Iu
      gJ7quVQKcfQcAVwuwWRrZXyhjhHaVKoPP5yRS+ESVTl70J5HBh6B7laooxf1yVAW
      8NJK1iQ1Pw2x3ABBo1cSMcTQ3Hk1ZWThJ7oPul2+QyzvOjIjiEPBstyzEPaxPG4I
      nH9ttalAwSLBsobVaK8mmiAdtAdx+CmHWrB4UNxCPYasrt5A6a9A9SiQ2dLd
      -----END CERTIFICATE-----
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("sets the CA cert field", func() {
						installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).ToNot(HaveOccurred())

						Expect(installationManifest.Cert.CA).To(Equal(`-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQLzf5Fs3v+Dblm+CKQFxiKTANBgkqhkiG9w0BAQsFADAm
MQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkwHhcNMTcwNTE2
MTUzNTI4WhcNMTgwNTE2MTUzNTI4WjAmMQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoT
DUNsb3VkIEZvdW5kcnkwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC+
4E0QJMOpQwbHACvrZ4FleP4/DMFvYUBySfKzDOgd99Nm8LdXuJcI1SYHJ3sV+mh0
+cQmRt8U2A/lw7bNU6JdM0fWHa/2nGjSBKWgPzba68NdsmwjqUjLatKpr1yvd384
PJJKC7NrxwvChgB8ui84T4SrXHCioYMDEDIqLGmHJHMKnzQ17nu7ECO4e6QuCfnH
RDs7dTjomTAiFuF4fh4SPgEDMGaCE5HZr4t3gvc9n4UftpcCpi+Jh+neRiWx+v37
ZAYf2kp3wWtYDlgWk06cZzHZZ9uYZFwHDNHdDKHxGGvAh2Rm6rpPF2oA6OEyx6BH
85/STCgSMCnV1Wkd+1yPAgMBAAGjIzAhMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMB
Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBGvGggx3IM4KCMpVDSv9zFKX4K
IuCRQ6VFab3sgnlelMFaMj3+8baJ/YMko8PP1wVfUviVgKuiZO8tqL00Yo4s1WKp
x3MLIG4eBX9pj0ZVRa3kpcF2Wvg6WhrzUzONf7pfuz/9avl77o4aSt4TwyCvM4Iu
gJ7quVQKcfQcAVwuwWRrZXyhjhHaVKoPP5yRS+ESVTl70J5HBh6B7laooxf1yVAW
8NJK1iQ1Pw2x3ABBo1cSMcTQ3Hk1ZWThJ7oPul2+QyzvOjIjiEPBstyzEPaxPG4I
nH9ttalAwSLBsobVaK8mmiAdtAdx+CmHWrB4UNxCPYasrt5A6a9A9SiQ2dLd
-----END CERTIFICATE-----
`))
					})
				})

				Context("that is invalid", func() {
					BeforeEach(func() {
						err := fakeFs.WriteFileString(comboManifestPath, `
---
name: fake-deployment-name
cloud_provider:
  template:
    name: fake-cpi-job-name
    release: fake-cpi-release-name
  mbus: http://fake-mbus-user:fake-mbus-password@0.0.0.0:6868
  cert:
    ca: |
      -----BEGIN CERTIFICATE-----
      no valid certificate
      -----END CERTIFICATE-----
`)
						Expect(err).ToNot(HaveOccurred())
						fakeUUIDGenerator.GeneratedUUID = "fake-uuid"
					})

					It("returns an error", func() {
						_, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(Equal("Invalid CA cert"))
					})
				})
			})

			Context("when ca cert is not provided", func() {
				BeforeEach(func() {
					err := fakeFs.WriteFileString(comboManifestPath, fixtures.missingPrivateKeyManifest)
					Expect(err).ToNot(HaveOccurred())
				})

				It("does not expand the path", func() {
					installationManifest, err := parser.Parse(comboManifestPath, boshtpl.StaticVariables{}, patch.Ops{}, releaseSetManifest)
					Expect(err).ToNot(HaveOccurred())

					Expect(installationManifest.Cert.CA).To(Equal(""))
				})
			})
		})
	})
})
