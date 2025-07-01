#!/bin/bash

# Test script for pod list functionality
# Creates a deployment, scales it up/down, then deletes it

DEPLOYMENT_NAME="test-deployment"
NAMESPACE="default"

echo "🚀 Starting deployment test..."
echo "This will create, scale, and delete a test deployment"
echo ""

# Function to wait and show status
wait_and_show() {
    echo "⏳ Waiting 2 seconds..."
    sleep 2
    echo "📊 Current pod status:"
    kubectl get pods -l app=$DEPLOYMENT_NAME --no-headers | wc -l | xargs echo "   Active pods:"
    echo ""
}

# Create deployment
echo "1️⃣ Creating deployment '$DEPLOYMENT_NAME' with 2 replicas..."
kubectl create deployment $DEPLOYMENT_NAME --image=nginx:latest --replicas=2
wait_and_show

# Scale up to 4
echo "2️⃣ Scaling up to 4 replicas..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=4
wait_and_show

# Scale up to 6
echo "3️⃣ Scaling up to 6 replicas..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=6
wait_and_show

# Scale down to 3
echo "4️⃣ Scaling down to 3 replicas..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=3
wait_and_show

# Scale down to 1
echo "5️⃣ Scaling down to 1 replica..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=1
wait_and_show

# Scale up to 5
echo "6️⃣ Scaling up to 5 replicas..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=5
wait_and_show

# Scale down to 2
echo "7️⃣ Scaling down to 2 replicas..."
kubectl scale deployment $DEPLOYMENT_NAME --replicas=2
wait_and_show

# Delete deployment
echo "8️⃣ Deleting deployment '$DEPLOYMENT_NAME'..."
kubectl delete deployment $DEPLOYMENT_NAME
wait_and_show

echo "✅ Test completed!"
echo "Check your vigilant-debug.log for detailed pod events" 