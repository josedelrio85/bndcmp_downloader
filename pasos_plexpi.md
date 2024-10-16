# Pasos en plexpi

## 1- crear variable de entorno

export BASE_FOLDER=/home/jose/Descargas/bndcmp_downloader

## 2- crear network

docker network create bndcmp_downloader_network -d bridge

## 3- crear volumen

docker volume create \
--driver local \
--opt type=none \
--opt device=$BASE_FOLDER \
--opt o=bind \
bndcmp_downloader_volume

## 4- build backend

docker build -t bndcmp_downloader_api -f backend/dockerfile ./backend

## 5 - run backend

docker run -d \
--name bndcmp_downloader_api \
-e BASE_FOLDER=$BASE_FOLDER \
-v bndcmp_downloader_volume:$BASE_FOLDER \
-p 8099:8099 \
--network bndcmp_downloader_network \
--restart unless-stopped \
bndcmp_downloader_api

## 6 - build frontend

docker build -t bndcmp_downloader_ui -f frontend/dockerfile ./frontend

## 7 - run frontend

docker run -d --name ui -p 8080:80 --network bndcmp_downloader_network --restart unless-stopped bndcmp_downloader_ui