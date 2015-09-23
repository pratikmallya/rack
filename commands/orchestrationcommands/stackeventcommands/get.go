package stackeventcommands

import (
	"github.com/rackspace/rack/commandoptions"
	"github.com/rackspace/rack/commands/orchestrationcommands/stackcommands"
	"github.com/rackspace/rack/handler"
	"github.com/rackspace/rack/internal/github.com/codegangsta/cli"
	"github.com/rackspace/rack/internal/github.com/rackspace/gophercloud/rackspace/orchestration/v1/stackevents"
	"github.com/rackspace/rack/util"
)

var get = cli.Command{
	Name:        "get",
	Usage:       util.Usage(commandPrefix, "get", "[--name <stackName> | --id <stackID>] --resource <resourceName> --event <eventID>"),
	Description: "Show details for a specified event",
	Action:      actionGet,
	Flags:       commandoptions.CommandFlags(flagsGet, keysGet),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsGet, keysGet))
	},
}

func flagsGet() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't specified] The stack name.",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't specified] The stack id.",
		},
		cli.StringFlag{
			Name:  "resource",
			Usage: "[required] The resource name.",
		},
		cli.StringFlag{
			Name:  "event",
			Usage: "[required] The event id.",
		},
	}
}

type paramsGet struct {
	stackName    string
	stackID      string
	resourceName string
	eventID      string
}

var keysGet = []string{"ResourceName", "Time", "ResourceStatusReason", "ResourceStatus", "PhysicalResourceID", "ID", "ResourceProperties"}

type commandGet handler.Command

func actionGet(c *cli.Context) {
	command := &commandGet{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandGet) Context() *handler.Context {
	return command.Ctx
}

func (command *commandGet) Keys() []string {
	return keysGet
}

func (command *commandGet) ServiceClientType() string {
	return serviceClientType
}

func (command *commandGet) HandleFlags(resource *handler.Resource) error {
	err := command.Ctx.CheckFlagsSet([]string{"resource", "event"})
	if err != nil {
		return err
	}
	c := command.Ctx.CLIContext
	name := c.String("name")
	id := c.String("id")
	name, id, err = stackcommands.IDAndName(command.Ctx.ServiceClient, name, id)
	if err != nil {
		return err
	}
	resource.Params = &paramsGet{
		stackName:    name,
		stackID:      id,
		resourceName: c.String("resource"),
		eventID:      c.String("event"),
	}
	return nil
}

func (command *commandGet) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsGet)
	stackName := params.stackName
	stackID := params.stackID
	resourceName := params.resourceName
	eventID := params.eventID

	result, err := stackevents.Get(command.Ctx.ServiceClient, stackName, stackID, resourceName, eventID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = eventSingle(result)
}
