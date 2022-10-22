package cloudfront

import (
    "log"
    "context"

    //"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

func GetDistributions(profile string) []string {
    // Load config based on a selected profile 
    cfg, err := config.LoadDefaultConfig(context.TODO(), 
    config.WithSharedConfigProfile(profile))

    if err != nil {
        log.Fatalf("Error loading config: %v", err)       
    }
    
    client := cloudfront.NewFromConfig(cfg)
    input := &cloudfront.ListDistributionsInput{}

    result, err := client.ListDistributions(context.Background(), input)

    return result

}
