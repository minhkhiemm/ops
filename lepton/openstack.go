package lepton

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// OpenStack provides access to the OpenStack API.
type OpenStack struct {
	Storage  *Datastores
	provider *gophercloud.ProviderClient
}

// ResizeImage is not supported on OpenStack.
func (o *OpenStack) ResizeImage(ctx *Context, imagename string, hbytes string) error {
	return fmt.Errorf("Operation not supported")
}

// BuildImage to be upload on OpenStack
func (o *OpenStack) BuildImage(ctx *Context) (string, error) {
	c := ctx.config
	err := BuildImage(*c)
	if err != nil {
		return "", err
	}

	return o.customizeImage(ctx)
}

// BuildImageWithPackage to upload on OpenStack.
func (o *OpenStack) BuildImageWithPackage(ctx *Context, pkgpath string) (string, error) {
	c := ctx.config
	err := BuildImageFromPackage(pkgpath, *c)
	if err != nil {
		return "", err
	}
	return o.customizeImage(ctx)
}

func (o *OpenStack) createImage(key string, bucket string, region string) {
	fmt.Println("creating image")

	// https://docs.openstack.org/api-ref/image/v2/?expanded=show-image-schema-detail#stage-binary-image-data

	// sounds like this is 3 steps:
	// 1) create a new image record
	// 2) stage it (upload it)
	// https://github.com/gophercloud/gophercloud/blob/6e3895ed427a63be0fe60dcb8a31e564495ee6e9/acceptance/openstack/imageservice/v2/imageservice.go#L116
	// 3) import it
	// https://github.com/gophercloud/gophercloud/blob/6e3895ed427a63be0fe60dcb8a31e564495ee6e9/acceptance/openstack/imageservice/v2/imageservice.go#L77

	imageClient, err := openstack.NewImageServiceV2(o.provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println(err)
	}

	visibility := images.ImageVisibilityPrivate

	createOpts := images.CreateOpts{
		Name:       "gtest",
		DiskFormat: "raw",
		Visibility: &visibility,
	}

	image, err := images.Create(imageClient, createOpts).Extract()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", image)
	fmt.Println("un-implemented")
}

// Initialize OpenStack related things
func (o *OpenStack) Initialize() error {

	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println(err)
	}

	o.provider, err = openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// CreateImage - Creates image on OpenStack using nanos images
// This merely uploads the flat and base image to the datastore and then
// creates a copy of the image to perform the vmfs translation (import
// does not do this by default). This sidesteps the vmfkstools
// transformation.
func (o *OpenStack) CreateImage(ctx *Context) error {
	fmt.Println("un-implemented")
	return nil
}

// ListImages lists images on a datastore.
// This is incredibly naive at the moment and probably worth putting
// under a root folder.
// essentially does the equivalent of 'govc datastore.ls'
func (o *OpenStack) ListImages(ctx *Context) error {

	imageClient, err := openstack.NewImageServiceV2(o.provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println(err)
	}

	listOpts := images.ListOpts{}

	allPages, err := images.List(imageClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Created"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})
	table.SetRowLine(true)

	for _, image := range allImages {
		var row []string
		row = append(row, image.Name)
		row = append(row, fmt.Sprintf("%v", image.Status))
		row = append(row, time2Human(image.CreatedAt))
		table.Append(row)
	}

	table.Render()

	return nil
}

// DeleteImage deletes image from OpenStack
func (o *OpenStack) DeleteImage(ctx *Context, imagename string) error {
	/*
		imageID := "1bea47ed-f6a9-463b-b423-14b9cca9ad27"
		err := images.Delete(imageClient, imageID).ExtractErr()
		if err != nil {
			panic(err)
		}
	*/

	fmt.Println("un-implemented")
	return nil
}

// CreateInstance - Creates instance on OpenStack.
// Currently we support pvsci adapter && vmnetx3 network driver.
func (o *OpenStack) CreateInstance(ctx *Context) error {
	client, err := openstack.NewComputeV2(o.provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println(err)
	}

	server, err := servers.Create(client, servers.CreateOpts{
		Name:      "My new server!",
		FlavorRef: "m1.micro",
		ImageRef:  "image_id",
	}).Extract()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", server)

	fmt.Println("un-implemented")

	return nil
}

// ListInstances lists instances on OpenStack.
// It essentially does:
// govc ls /ha-datacenter/vm
func (o *OpenStack) ListInstances(ctx *Context) error {
	fmt.Println("un-implemented")
	return nil
}

// DeleteInstance deletes instance from OpenStack
func (o *OpenStack) DeleteInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")
	return nil
}

// StartInstance starts an instance in OpenStack.
// It is the equivalent of:
// govc vm.power -on=true <instance_name>
func (o *OpenStack) StartInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")
	return nil
}

// StopInstance stops an instance from OpenStack
// It is the equivalent of:
// govc vm.power -on=false <instance_name>
func (o *OpenStack) StopInstance(ctx *Context, instancename string) error {
	fmt.Println("un-implemented")
	return nil

}

// GetInstanceLogs gets instance related logs.
// govc datastore.tail -n 100 gtest/serial.out
// logs don't appear until you spin up the instance.
func (o *OpenStack) GetInstanceLogs(ctx *Context, instancename string, watch bool) error {
	fmt.Println("un-implemented")
	return nil

}

// Todo - make me shared
func (o *OpenStack) customizeImage(ctx *Context) (string, error) {
	imagePath := ctx.config.RunConfig.Imagename
	return imagePath, nil
}
