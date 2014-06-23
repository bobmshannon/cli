package appsecuritygroup_test

import (
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	testapi "github.com/cloudfoundry/cli/testhelpers/api"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/appsecuritygroup"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("application-security-group command", func() {
	var (
		ui                   *testterm.FakeUI
		appSecurityGroupRepo *testapi.FakeAppSecurityGroup
		requirementsFactory  *testreq.FakeReqFactory
		configRepo           configuration.ReadWriter
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		appSecurityGroupRepo = &testapi.FakeAppSecurityGroup{}
		configRepo = testconfig.NewRepositoryWithDefaults()
	})

	runCommand := func(args ...string) {
		cmd := NewShowAppSecurityGroup(ui, configRepo, appSecurityGroupRepo)
		testcmd.RunCommand(cmd, args, requirementsFactory)
	}

	Describe("requirements", func() {
		It("should fail if not logged in", func() {
			runCommand("my-group")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

		It("should fail with usage when not provided a single argument", func() {
			requirementsFactory.LoginSuccess = true
			runCommand("whoops", "I can't believe", "I accidentally", "the whole thing")
			Expect(ui.FailedWithUsage).To(BeTrue())
		})
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		Context("when the group with the given name exists", func() {
			BeforeEach(func() {
				rulesMap := []map[string]string{{"just-pretend": "that-this-is-correct"}}

				appSecurityGroupRepo.ReadReturns.ApplicationSecurityGroup = models.ApplicationSecurityGroup{
					ApplicationSecurityGroupFields: models.ApplicationSecurityGroupFields{
						Name:  "my-group",
						Guid:  "group-guid",
						Rules: rulesMap,
					},
					Spaces: []models.SpaceFields{{Name: "space-1"}, {Name: "space-2"}},
				}
			})

			It("should fetch the application security group from its repo", func() {
				runCommand("my-group")
				Expect(appSecurityGroupRepo.ReadCalledWith.Name).To(Equal("my-group"))
			})

			It("tells the user what it's about to do and then shows the group", func() {
				runCommand("my-group")
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Getting", "application security group", "my-group", "my-user"},
					[]string{"OK"},
					[]string{"Name:", "my-group"},
					[]string{"Rules:", `[{"just-pretend":"that-this-is-correct"}]`},
					[]string{"Spaces:", "space-1, space-2"},
				))
			})

			It("tells the user if no spaces are assigned", func() {
				appSecurityGroupRepo.ReadReturns.ApplicationSecurityGroup = models.ApplicationSecurityGroup{
					ApplicationSecurityGroupFields: models.ApplicationSecurityGroupFields{
						Name:  "my-group",
						Guid:  "group-guid",
						Rules: []map[string]string{},
					},
					Spaces: []models.SpaceFields{},
				}

				runCommand("my-group")

				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Spaces:", "No spaces"},
				))
			})
		})

		It("fails and warns the user if a group with that name could not be found", func() {
			appSecurityGroupRepo.ReadReturns.Error = errors.New("half-past-tea-time")
			runCommand("im-late!")

			Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
		})
	})
})