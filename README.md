# GinChat

![GinChat](https://img.shields.io/badge/Go-1.22-blue)
![Docker](https://img.shields.io/badge/Docker-Container-green)
![Redis](https://img.shields.io/badge/Redis-7.0.11-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16.2-orange)

> A real-time chat application built using the Gin framework (Go), PostgreSQL, and Redis.

GinChat is a web chat and efficient real-time chat service leveraging the Go Gin web framework, Redis for message caching, and PostgreSQL as the database. The application is containerized using Docker for easy deployment.

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup and Installation](#setup-and-installation)
    - [Prerequisites](#prerequisites)
    - [Run Using Docker Compose](#run-using-docker-compose)
    - [Manual Setup](#manual-setup)
- [Environment Variables](#environment-variables)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

---

## Features
- ❤️ **Technologies in Use**: Gorilla WebSocket, GORM, Swagger, concurrency, JWT, OTP, microservice architecture.
- 💬 **Chat Features**: Safe and private chats where you can send text, images, and voice messages.
- 🟢 **Real-Time Communication**: Enables instant messaging with Redis.
- ⚡ **Gin Framework**: Lightweight and high-performance Go framework.
- 🔐 **Authentication**: Secure user authentication with JWT and OTP via phone numbers.
- 🛠 **Dockerized**: All services (GinChat, PostgreSQL, Redis) run in isolated containers.
- 💾 **Persistence**: PostgreSQL for persistent message storage and Redis for in-memory caching to reduce load on PostgreSQL.
- 📊 **Scalable**: Easily scalable with Docker and Redis.

---

## Tech Stack
- **Language**: Go (Gin framework)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Containerization**: Docker
- **Frontend**: Currently, there is no UI.

---

## Setup and Installation

### Prerequisites
- **Go** (v1.22+)
- **Docker** (Latest version)
- **Docker Compose** (v2+)

### Run Using Docker Compose
1. **Clone the repository:**
   ```bash
   git clone https://github.com/Ali-darzii/GinChat.git
2. **Go to the main directory :**
   ```bash
   cd GinChat
3. **Run with Docker Compose :**
   ```bash
   docker compose up -d --build
4. **If you have an older version of Docker Compose :**
    ```bash
   docker-compose up -d --build

5. **For API Endpoints check :**
   [Swagger Documentation](http://127.0.0.1:8080/swagger/index.html)
