package cloudfront

import (
    "log"
    "context"

    "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

func GetDistributions(profile string) []string {
    // Load config based on a selected profile 
    cfg, err := config.LoadDefaultConfig(context.TODO(), 
    config.WithSharedConfigProfile(profile))

    if err != nil {
        log.Fatalf("Error loading config: %v", err)       
    }
    
    // Create a new session to use with CloudFront service
    s, err := session.NewSession(cfg)
    if err != nil {
        log.Fatal(err)
    }

    client := cloudfront.New(s)
    input := &cloudfront.ListDistributionsInput{}

    result, err := 


}
