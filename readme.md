# Seleksi Tim Laboratorium Programming 2019 Tahap II
> Implementasi CRUD website untuk mengelola dorayaki dan toko

## Table of contents
* [General info](#general-info)
* [Swagger](#swagger)
* [Technologies](#technologies)
* [Setup](#setup)
* [Features](#features)
* [Status](#status)
* [Inspiration](#inspiration)
* [Contact](#contact)

## General info
Project ini dibagi menjadi dua bagian, Backend dan Frontend. Project ini dibuat untuk memecahkan permasalahan manajemen konten yang diperlukan toko "Stand with Dorayaki". Project ini dibuat untuk memenuhi salah satu syarat dalam seleksi tim labpro 2019 yang kedua.
Repository ini untuk bagian backend.

## Swagger
Klik link berikut [ini](https://api.dorayaki.wiwid.me/docs/api/v1/)

## Technologies
* Docker
* Golang
## Setup
1. Install docker dan docker-compose di komputer.
Klik link [ini](https://docs.docker.com/engine/install/) untuk panduan menginstall docker
Klik link [ini](https://docs.docker.com/compose/install/) untuk panduan menginstall docker-compose.
2. Dengan docker sebenarnya sudah cukup karena golang akan diinstal di dalamnya. Jika ingin menginstall golang di host komputer Anda, cukup liat link [ini](https://golang.org/doc/install) untuk panduang menginstall golang
3. Lakukan git clone repository ini dengan mengetikkan di terminal atau git bash
```bash
git clone https://github.com/widyaput/dorayaki-backend.git
```

## Code Examples
Show examples of usage:
1. Masuk ke folder hasil clone git repository
2. Jalankan perintah ini di terminal
```bash
docker-compose -f deployments/compose/docker-compose.yml -p dorayakidev up --build
```
Jika pertama kali membuild docker, untuk aksi selanjutnya hilangkan tag --build sehingga hanya tersisa
```bash
docker-compose -f deployments/compose/docker-compose.yml -p dorayakidev up
```
saja

3. Untuk menjalankan seeder, Anda harus melakukan perintah
```bash
docker exec -it ${id container app} bash && go run cmd/server/main.go seed
```
id container dapat dilihat menggunakan perintah
```bash
docker ps
```
lalu pilih container yang sesuai

4. Bagian backend sudah siap dijalankan, defaultnya pada port 8080.
5. Anda dapat melihat swagger dalam path '/docs/api/v1'

## Features
List of features ready and TODOs for future development
* Dideploy pada [url](https://api.dorayaki.wiwid.me/) walau terkadang sering down apalagi ketika membuild bagian frontend. Disetting menggunakan Github Actions.
* Fitur search sekaligus paginasi pada dorayaki dan toko.
* Fitur login layaknya CMS pada umumnya. Silahkan cari credentialnya di package seeds biar seru :D

## Status
Project is: _finished_

## Contact
Created by Widya Anugrah Putra - feel free to contact me!