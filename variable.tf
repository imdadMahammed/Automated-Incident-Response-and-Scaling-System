variable "region" {
  description = "AWS region"
  default     = "us-west-2"
}

variable "vpc_name" {
  description = "Name of the VPC"
  default     = "eks-vpc"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones for the VPC"
  type        = list(string)
  default     = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "public_subnets" {
  description = "Public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "private_subnets" {
  description = "Private subnets"
  type        = list(string)
  default     = ["10.0.4.0/24", "10.0.5.0/24", "10.0.6.0/24"]
}

variable "cluster_name" {
  description = "Name of the EKS cluster"
  default     = "my-eks-cluster"
}

variable "key_name" {
  description = "Key name for SSH access to the EKS nodes"
  default     = "my-eks-key"
}

variable "s3_bucket_name" {
  description = "S3 bucket name for Prometheus storage"
  default     = "prometheus-backups-bucket"
}

variable "slack_webhook_url" {
  description = "Slack webhook URL for alert notifications"
  type        = string
}
