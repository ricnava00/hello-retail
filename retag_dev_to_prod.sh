#! /bin/bash

PROD_TAG=accesscontrol

components='gateway'

for component in $components
do
	old_tag="openfaas/${component}:latest-dev"
	new_tag="openfaasaccesscontrol/${component}:latest-${PROD_TAG}"
	echo Retagging $old_tag to $new_tag
	docker tag $old_tag $new_tag
done


services='product-catalog-api product-catalog-builder product-photos-assign product-photos-message product-photos-receive product-photos-record product-photos-report product-photos-success product-purchase-authenticate product-purchase-authorize-cc product-purchase-get-price product-purchase-publish product-photos product-purchase'

for service in $services
do
	old_tag="${service}:latest"
	new_tag="openfaasaccesscontrol/${service}:latest-${PROD_TAG}"
	echo Retagging $old_tag to $new_tag
	docker tag $old_tag $new_tag
done
