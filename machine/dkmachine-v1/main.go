package main

import (
	"fmt"
	"os"

	"github.com/otiai10/dkmachine"
	"github.com/otiai10/jsonindent"
)

func main() {
	opt := &dkmachine.CreateOptions{
		Name:                        "dkmachine-test",
		Driver:                      "amazonec2",
		AmazonEC2Region:             "ap-northeast-1",
		AmazonEC2RootSize:           8,
		AmazonEC2InstanceType:       "t2.nano",
		AmazonEC2SecurityGroup:      "awsub-default",
		AmazonEC2IAMInstanceProfile: "awsubtest",
	}

	m, err := dkmachine.Create(opt)
	if err != nil {
		panic(err)
	}

	jsonindent.NewEncoder(os.Stdout).Encode(m)
	// defer m.Remove()

	// fmt.Printf("%+v\n", m.HostConfig.Driver)

	// jsonindent.NewEncoder(os.Stdout).Encode(m.HostConfig.Driver)

	fmt.Println(
		m.GetPrivateIPAddress(),
	)
}
