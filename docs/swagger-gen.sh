#!/bin/sh

cd "$(dirname "$0")"
DOMAIN_FILES=`find ../internal/domain -name "*.yaml" \! -name common.yaml`
MODULE=`cat ../go.mod | grep module | sed 's/^module //'`
SWAGGER_FILES="swagger.yaml ../internal/domain/common/common.yaml"

# ensure required binaries are installed
command -v oapi-codegen && echo "oapi-codegen installed, nothing to do" || ( echo "install oapi-codegen"; go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.1)
command -v yq  && echo "yq installed, nothing to do" || (echo "install yq"; go install github.com/mikefarah/yq/v2@latest)

# generate model & handler for each domain
for FILE in ${DOMAIN_FILES}; do
	DOMAIN=`echo ${FILE} | sed -e 's/\.\.\/internal\/domain\///' -e 's/\/.*//'`
	echo "\nparse: ${DOMAIN}.yaml -> ${DOMAIN}_gen.go"

	if [ -f ../internal/domain/${DOMAIN}/${DOMAIN}-readonly.yaml ]; then
		SWAGGER_FILES="${SWAGGER_FILES} ../internal/domain/${DOMAIN}/${DOMAIN}-readonly.yaml"
	else
		${GOPATH}/bin/oapi-codegen \
		-package=${DOMAIN} \
		-import-mapping="../common/common.yaml:${MODULE}/internal/domain/common" \
		-generate='types,server' \
		../internal/domain/${DOMAIN}/${DOMAIN}.yaml \
		> ../internal/domain/${DOMAIN}/${DOMAIN}_gen.go

		SWAGGER_FILES="${SWAGGER_FILES} ../internal/domain/${DOMAIN}/${DOMAIN}.yaml"
	fi
done

# assemble swagger.gen.yaml
echo "\nassemble swagger_gen.yaml"
${GOPATH}/bin/yq merge ${SWAGGER_FILES} > ../internal/service/server/swagger/swagger_gen.yaml

# overwrite placeholder
echo "overwrite placeholder in swagger_gen.yaml"
cat ../internal/service/server/swagger/swagger_gen.yaml \
	| sed -e "s/__VERSION__/${VERSION}/" -e "s/__APP_NAME__/${APP_NAME}/" \
	| sed -e "s/..\/common\/common.yaml#\(.*\)/'#\1'/" \
	> out.yaml
mv out.yaml	../internal/service/server/swagger/swagger_gen.yaml
