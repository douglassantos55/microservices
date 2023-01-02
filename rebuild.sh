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

    return 0
}

rebuild_service
