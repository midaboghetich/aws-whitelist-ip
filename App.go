package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

var (
	flgVersion bool
)

func main() {

	const defaultSg = "sg-000000000"
	const defaultEmail = "undefined@mail.com"
	const defaultIpAddress = "ipAddress"

	var portArr flagArray

	secgroup := flag.String("g", defaultSg, "the security group where the ip should be whitelisted (mandatory)")
	email := flag.String("email", defaultEmail, "the email address used to identify the whitelisted IP Address (mandatory)")
	ipAddress := flag.String("ip", defaultIpAddress, "the ip address to whitelist")

	flag.Var(&portArr, "p", "port number to whitelist (mandatory)")
	flag.BoolVar(&flgVersion, "v", false, "if true, print version and exit")

	flag.Parse()

	if flgVersion {
		fmt.Printf("Build version %s %s\n", BuildVersion, BuildTime)
		os.Exit(0)
	}

	if *secgroup == defaultSg {
		flag.Usage()
		log.Fatal("you must specify the security group Id")
	}

	if *email == defaultEmail {
		flag.Usage()
		log.Fatal("you must specify your email address")
	}

	if len(portArr) == 0 {
		flag.Usage()
		log.Fatal("you must specify at least one port number")
	}

	mySecGroup := *secgroup
	description := *email

	var ip string

	if *ipAddress == defaultIpAddress {
		ip = discoverIp()
	} else {
		ip = *ipAddress
	}

	protocol := "TCP"

	cidrIpToWhitelist := ip + "/32"

	cfg, err := external.LoadDefaultAWSConfig()

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	ec2Svc := ec2.New(cfg)

	revokeAllIngressIP(ec2Svc, mySecGroup, description, protocol)

	for i :=0  ;i < len(portArr) ; i++ {
		tmpPort := portArr[i]

		currPort, err := strconv.ParseInt(tmpPort, 10, 64)
		if err != nil {
			panic("Oops, port value, invalid parameter")
		}

		addSecurityGroupIngress(ec2Svc, mySecGroup, protocol, currPort, currPort, cidrIpToWhitelist, description)
	}


	fmt.Println("Successfully set security group ingress")
}

func addSecurityGroupIngress(ec2Svc *ec2.EC2, mySecGroup string, protocol string, fromPort int64, toPort int64, cidrIpToWhitelist string, description string) {
	ingressRequest := ec2Svc.AuthorizeSecurityGroupIngressRequest(&ec2.AuthorizeSecurityGroupIngressInput{

		GroupId: &mySecGroup,
		IpPermissions: []ec2.IpPermission{
			{
				IpProtocol: &protocol,
				FromPort:   &fromPort,
				ToPort:     &toPort,
				IpRanges: []ec2.IpRange{
					{
						CidrIp:      &cidrIpToWhitelist,
						Description: &description,
					},
				},
			},
		},
	})
	_, err := ingressRequest.Send()

	if err != nil {
		log.Fatal("Unable to set security group %q ingress, %v", mySecGroup, err)
	}
}

func revokeAllIngressIP(ec2Svc *ec2.EC2, mySecGroup string, description string, protocol string) {

	request := ec2Svc.DescribeSecurityGroupsRequest(&ec2.DescribeSecurityGroupsInput{GroupIds: []string{mySecGroup}})

	output, err := request.Send()

	if err != nil {
		panic("something went terribly wrong, " + err.Error())

	}

	for _, group := range output.SecurityGroups {
		for _, permission := range group.IpPermissions {
			deleteFromPort := permission.FromPort
			deleteToPort := permission.ToPort

			for _, ipRange := range permission.IpRanges {
				if strings.Contains(*ipRange.Description, description) {

					fmt.Println("Revoking ip access")
					fmt.Println(ipRange)
					revokeIngressRequest := ec2Svc.RevokeSecurityGroupIngressRequest(&ec2.RevokeSecurityGroupIngressInput{

						GroupId: &mySecGroup,
						IpPermissions: []ec2.IpPermission{
							{
								IpProtocol: &protocol,
								FromPort:   deleteFromPort,
								ToPort:     deleteToPort,
								IpRanges: []ec2.IpRange{
									{
										CidrIp:      ipRange.CidrIp,
										Description: &description,
									},
								},
							},
						},
					})

					_, err := revokeIngressRequest.Send()

					if err != nil {
						log.Fatal("Unable to revoke IP address: %v", err)
					}

					fmt.Println("Successfully revoke security group ingress")

				}
			}
		}

	}
}


