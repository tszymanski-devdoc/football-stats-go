# Cloud Run Deployment Guide

## Prerequisites

Before deploying to Google Cloud Run, ensure you have:

1. **GCP Project** with billing enabled
2. **GitHub Repository Secrets** configured:
   - `GCP_SA_KEY` - Service Account JSON key
   - `GCP_PROJECT_ID` - Your GCP Project ID
   - `GCP_REGION` - Deployment region (optional, defaults to us-central1)

## Setting Up GitHub Secrets

### 1. Create a GCP Service Account

```bash
# Set your project ID
export PROJECT_ID="your-project-id"

# Create service account
gcloud iam service-accounts create github-actions \
    --display-name="GitHub Actions Deployer"

# Grant necessary permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"

# Create and download key
gcloud iam service-accounts keys create key.json \
    --iam-account=github-actions@${PROJECT_ID}.iam.gserviceaccount.com
```

### 2. Add Secrets to GitHub

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret** and add:
   - **Name**: `GCP_SA_KEY`
   - **Value**: Contents of `key.json` file (entire JSON)
   - **Name**: `GCP_PROJECT_ID`
   - **Value**: Your GCP project ID
   - **Name**: `GCP_REGION` (optional)
   - **Value**: `us-central1` (or your preferred region)

### 3. Enable Required GCP APIs

```bash
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable cloudbuild.googleapis.com
```

## Deployment Workflow

The deployment happens automatically through GitHub Actions:

### CI Pipeline (api-ci.yml)
Triggered on push/PR to `main` affecting `football-stats-go-api/**`:
1. ✅ Checkout code
2. ✅ Set up Go 1.24
3. ✅ Download dependencies
4. ✅ Run linting (go vet)
5. ✅ Build application
6. ✅ Run tests with coverage
7. ✅ Build Docker image
8. ✅ Push to Google Container Registry

### CD Pipeline (api-cd.yml)
Triggered after successful CI pipeline:
1. ✅ Deploy to Cloud Run
2. ✅ Configure auto-scaling (0-10 instances)
3. ✅ Set resource limits (512Mi memory, 1 CPU)
4. ✅ Test health endpoint
5. ✅ Generate deployment summary

## Manual Deployment

If you need to deploy manually:

```bash
# Build the image
cd football-stats-go-api
gcloud builds submit --tag gcr.io/$PROJECT_ID/football-stats-go-api

# Deploy to Cloud Run
gcloud run deploy football-stats-go-api \
  --image gcr.io/$PROJECT_ID/football-stats-go-api:latest \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --memory 512Mi \
  --cpu 1 \
  --max-instances 10
```

## Deployment Features

✅ **Auto-scaling**: 0 to 10 instances based on traffic  
✅ **Resource limits**: 512Mi memory, 1 CPU per instance  
✅ **Health checks**: Automatic verification post-deployment  
✅ **Zero-downtime**: Rolling updates  
✅ **HTTPS**: Automatic SSL certificates  
✅ **Public access**: No authentication required  

## Monitoring

After deployment, monitor your service:

```bash
# View logs
gcloud run logs tail football-stats-go-api --region us-central1

# Get service details
gcloud run services describe football-stats-go-api --region us-central1

# List revisions
gcloud run revisions list --service football-stats-go-api --region us-central1
```

## Cost Optimization

Cloud Run pricing:
- **Free tier**: 2 million requests/month
- **Pricing**: Pay only when processing requests
- **Auto-scaling to zero**: No cost when idle

## Troubleshooting

### Build Fails
- Check `go.mod` and `go.sum` are committed
- Verify Go version is 1.24
- Ensure `cmd/api/main.go` exists

### Deployment Fails
- Verify service account has correct permissions
- Check GCP APIs are enabled
- Ensure `GCP_SA_KEY` secret is valid JSON
- Verify image was pushed to GCR

### Service Not Responding
- Check logs: `gcloud run logs tail football-stats-go-api`
- Verify port 8080 is exposed in Dockerfile
- Test health endpoint: `curl https://your-service-url/health`

## Next Steps

1. **Add Custom Domain**: Configure Cloud Run with your domain
2. **Add Secrets**: Use Secret Manager for sensitive data
3. **Add Database**: Connect to Cloud SQL
4. **Add Monitoring**: Set up Cloud Monitoring alerts
5. **Add CI/CD for staging**: Create separate workflows for dev/staging/prod

## Useful Commands

```bash
# Get service URL
gcloud run services describe football-stats-go-api \
  --region us-central1 \
  --format 'value(status.url)'

# Update environment variables
gcloud run services update football-stats-go-api \
  --region us-central1 \
  --set-env-vars "NEW_VAR=value"

# Scale instances
gcloud run services update football-stats-go-api \
  --region us-central1 \
  --min-instances 1 \
  --max-instances 20

# Delete service
gcloud run services delete football-stats-go-api --region us-central1
```
