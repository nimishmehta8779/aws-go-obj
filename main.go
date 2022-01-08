package main

import (
	"github.com/nimishmehta8779/aws-go-obj/vpc"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		_, err := vpc.NewVpc(ctx, "myvpc", &vpc.VpcInput{
			VpcCidrBlock:           "10.10.0.0/16",
			AvailabilityZone:       []string{"us-east-1a", "us-east-1b"},
			PrivateSubnetCidrBlock: []string{"10.10.20.0/24", "10.10.30.0/24"},
		})
		if err != nil {
			return err
		}
		return nil
	},
	)
}
