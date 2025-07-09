---
name: deploy
description: Build, push, and deploy my-awesome-app to the target environment
---

# Deploy

Automates container build, image publish, and service update for **my-awesome-app**.

## Commands

```bash
# Change to project root
cd $(git rev-parse --show-toplevel)

# Build Docker image
docker build -t my-awesome-app:${CI_COMMIT_SHA:-local} .

# Authenticate to registry
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789012.dkr.ecr.us-east-1.amazonaws.com

# Tag and push
docker tag my-awesome-app:${CI_COMMIT_SHA:-local} 123456789012.dkr.ecr.us-east-1.amazonaws.com/my-awesome-app:${CI_COMMIT_SHA:-local}
docker push 123456789012.dkr.ecr.us-east-1.amazonaws.com/my-awesome-app:${CI_COMMIT_SHA:-local}

# Deploy to ECS
aws ecs update-service \
  --cluster my-awesome-cluster \
  --service my-awesome-service \
  --force-new-deployment \
  --region us-east-1
```

## Configuration

- Ensure AWS CLI is configured with appropriate credentials.
- ECR repository `123456789012.dkr.ecr.us-east-1.amazonaws.com/my-awesome-app` exists.
- ECS cluster `my-awesome-cluster` and service `my-awesome-service` are defined in AWS.

## Troubleshooting

- **Permission denied**: verify AWS IAM permissions for ECR and ECS.
- **Image not found**: check that the image tag matches pushed tag.
- **Service unchanged**: ensure `--force-new-deployment` flag is present.