# automation-as-a-service
Day One Provisioning Automation

# Roadmap
1. Landing zone with all best practices - security, observability, compliance, HA
  - multiAZ networking
  - Cloudtrail and logs for it
  - VPC flow logs
  - Unified tagging
  - Cost allocation tagging
  - Compliance scanning - security hub
  - AWS Config compliance rules
  - VPC Endpoints - Gateway and Interface

2. AWS Organizations with 3 accounts
  - management
  - audit
  - application
  - SSO for login to all accounts
  - ability to extend with multiple application accounts that join AWS organization

3. Catalog of resources that can be provisioned, all following best practices
   - EKS
   - ECS
   - Cloudfront
   - Opensearch
   - RDS
   - DynamoDB
   - S3
   - R53
   - Elasticache
  
4. More rare and obscure services
   - Redshift
   - Data Streaming and analytics - Kinesis
   - ML
   - AI
  
5. Self service catalog frontend
