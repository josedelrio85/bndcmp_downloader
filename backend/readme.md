# Bandcamp Downloader

This project is a simple API that allows you to download music from Bandcamp. It uses a combination of web scraping and direct download links to get the music.

## Getting Started

# set the device path to the correct downloads folder in the host machine
docker volume create \
--driver local \
--opt type=none \
--opt device=/your/downloads/folder \
--opt o=bind \
bndcmp_downloader_volume

docker build -t bndcmp_downloader_api .

# set the container path [/app/downloads] matching the downloads folder name used in .env file

docker run \
--name bndcmp_downloader_api \
--env-file .env \
-v bndcmp_downloader_volume:/app/downloads \
-p 8099:8099 \
bndcmp_downloader_api

### Prerequisites


