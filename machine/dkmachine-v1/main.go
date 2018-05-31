package main

import (
	"fmt"

	"github.com/otiai10/dkmachine"
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

	fmt.Printf("%+v\n", m)
	// defer m.Remove()
}
