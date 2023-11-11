variable "s3_bucket" {
  type        = string
  description = "S3 bucket for storing lambda packages before deploy"
}

variable "s3_bucket_key" {
  type        = string
  description = "S3 bucket key with distrib archive"
  default     = ""
}

variable "name_prefix" {
  type        = string
  description = "Lambda function names prefix"
  default     = "terraform-registry"
}

variable "create_domain" {
  type    = bool
  default = false
}
variable "record_name" {
  type = string
}
variable "domain_name" {
  type = string
}

variable "access_log_format" {
  type    = string
  default = "$context.identity.sourceIp - - [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage"
}

variable "storage_bucket_name" {
  type = string
}

variable "storage_bucket_readonly_arns" {
  type = list(string)
}

variable "storage_bucket_readwrite_arns" {
  type = list(string)
}

variable "storage_bucket_admin_arns" {
  type = list(string)
}

# variable "kms_sign_alias" {
#   type = string
# }

# variable "kms_sign_user_arns" {
#   type = list(string)
# }

# variable "kms_sign_admin_arns" {
#   type = list(string)
# }

variable "users" {
  type = map(map(string))
}

# variable "openpgp_name_prefix" {
#   type = string
# }

# variable "openpgp_email" {
#   type = string
# }

variable "tags" {
  type        = map(any)
  description = "Additional tags to add to resources"
  default     = {}
}
