echo "Initiating infra setup for Zyg services..."
echo "setting up environment variables..."
export PROJECT_ID=$(gcloud config get-value project)
export PROJECT_NAME=$(gcloud config get-value project)
export PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')
export LOCATION=us-west1

IMAGE_REPO=backend

IMAGE_NAME=$LOCATION-docker.pkg.dev/$PROJECT_ID/$IMAGE_REPO/zyg-xsrv-image

docker tag zyg-xsrv-image $LOCATION-docker.pkg.dev/$PROJECT_ID/$IMAGE_REPO/zyg-xsrv-image:latest

docker build -t $IMAGE_NAME -f xsrv.DockerFile .

docker push $IMAGE_NAME

ENV_VARS_FILE=.env.yaml

gcloud run deploy zyg-backend-xsrv \
    --image $IMAGE_NAME \
    --platform managed \
    --region $LOCATION \
    --allow-unauthenticated \
    --env-vars-file $ENV_VARS_FILE \
