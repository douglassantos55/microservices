#!/usr/bin/bash

SERVICE=$1

rebuild_service()
{
    if [ "$SERVICE" == "" ]; then
        echo "missing service name"
        return -1
    fi

    minikube kubectl -- delete -f $SERVICE.yml
    minikube image build ./$SERVICE -t $SERVICE
    minikube kubectl -- apply -f $SERVICE.yml

    # clean up so that it doesn't pile up and fills the disk
    minikube ssh -- docker system prune --volumes -f

    return 0
}

rebuild_service
