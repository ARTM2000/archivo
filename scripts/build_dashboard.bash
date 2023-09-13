#!/bin/bash

install_dependencies() {
    printf "installing dependencies by node version '$(node -v)'\n\n"
    npm ci 
}

build_dashboard() {
    printf "building dashboard by base url '${1}'\n"
    VITE_ARCHIVO_API_PANEL_BASE_URL=${1} npm run build
}

cd web && install_dependencies && build_dashboard ${1}