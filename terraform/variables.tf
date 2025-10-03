variable "project_id" {
  description = "The GCP project ID"
  type        = string
  default     = "pluma-development"
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "europe-west1"
}

variable "zone" {
  description = "The GCP zone"
  type        = string
  default     = "europe-west1-b"
}

variable "cluster_name" {
  description = "The name of the GKE cluster"
  type        = string
  default     = "dev-pluma-cluster"
}

variable "namespace" {
  description = "Kubernetes namespace for the application"
  type        = string
  default     = "pluma-kanban"
}

variable "artifact_registry_repository" {
  description = "Artifact Registry repository name"
  type        = string
  default     = "project-t-repo"
}
