package vpc

import (
	"fmt"

	util "github.com/nimishmehta8779/aws-go-obj/utils"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type VpcInput struct {
	VpcCidrBlock           string
	AvailabilityZone       []string
	PrivateSubnetCidrBlock []string
}

func (in *VpcInput) Validate() error {
	azCount := len(in.AvailabilityZone)
	if len(in.PrivateSubnetCidrBlock) > azCount {
		return fmt.Errorf("not Enough availability zones are provided")
	}
	return nil
}

type VpcOutput struct {
	Vpc            *ec2.Vpc
	PrivateSubnets []*ec2.Subnet
	Subnet0        string
	Subnet1        string
}

func NewVpc(ctx *pulumi.Context, name string, input *VpcInput) (*VpcOutput, error) {
	var err error
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("while validating input %v", err)
	}

	output := &VpcOutput{}

	output.Vpc, err = ec2.NewVpc(ctx, name, &ec2.VpcArgs{
		CidrBlock: pulumi.String(input.VpcCidrBlock),
		Tags:      pulumi.ToStringMap(util.NewNameTags(ctx, name)),
	})

	if err != nil {
		return nil, err
	}

	if input.PrivateSubnetCidrBlock != nil {
		if err := NewPrivateSubnet(ctx, input, output); err != nil {
			return nil, fmt.Errorf("while creating private subnet %v", err)
		}
	}
	ctx.Export("VPC-ID", output.Vpc.ID())
	return output, nil
}

func NewPrivateSubnet(ctx *pulumi.Context, input *VpcInput, output *VpcOutput) error {
	routes := make(ec2.RouteTableRouteArray, 0)

	// routes = append(routes, ec2.RouteTableRouteArgs{
	// 	CidrBlock: pulumi.String("0.0.0.0/0"),
	// })

	rt, err := ec2.NewRouteTable(ctx, "private-rt", &ec2.RouteTableArgs{
		VpcId:  output.Vpc.ID(),
		Routes: routes,
		Tags:   pulumi.ToStringMap(util.NewNameTags(ctx, "private-rt")),
	})
	if err != nil {
		return err
	}
	output.PrivateSubnets = make([]*ec2.Subnet, 0, len(input.PrivateSubnetCidrBlock))

	for i, cidr := range input.PrivateSubnetCidrBlock {
		az := input.AvailabilityZone[i]
		name := fmt.Sprintf("%s-%d", "private-subnet", i)
		subnet, err := ec2.NewSubnet(ctx, name, &ec2.SubnetArgs{
			VpcId:            output.Vpc.ID(),
			CidrBlock:        pulumi.String(cidr),
			AvailabilityZone: pulumi.String(az),
			Tags:             pulumi.ToStringMap(util.NewNameTags(ctx, name)),
		})
		if err != nil {
			return err
		}

		if _, err := ec2.NewRouteTableAssociation(ctx, name, &ec2.RouteTableAssociationArgs{
			RouteTableId: rt.ID(),
			SubnetId:     subnet.ID(),
		}); err != nil {
			return err
		}
		output.PrivateSubnets = append(output.PrivateSubnets, subnet)
	}
	privateSubnetIDs := make([]interface{}, 0)
	for _, subnet := range output.PrivateSubnets {
		privateSubnetIDs = append(privateSubnetIDs, subnet.ID().ToStringOutput())
	}
	ctx.Export("Subnet-ID", pulumi.All(privateSubnetIDs...).ApplyT(util.StringArrayOutputFunc))
	ctx.Export("Subnets", pulumi.All(privateSubnetIDs...).ApplyT(func(args []interface{}) string {
		return fmt.Sprintf("%s, %s", args[0], args[1])
	}))
	return nil
}
