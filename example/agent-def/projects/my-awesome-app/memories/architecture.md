---
name: architecture
description: Project architecture and context for my-awesome-app
---

# Project Architecture

## Overview
The **my-awesome-app** is a full-stack web application designed to manage user tasks and collaboration in real time. It follows a service-oriented architecture to separate concerns between frontend, backend, and infrastructure.

## Tech Stack
- Frontend: React (TypeScript), Redux Toolkit, Tailwind CSS  
- Backend: Go, Gin HTTP framework, GORM for ORM  
- Database: PostgreSQL  
- Real-time: WebSockets (Gorilla WebSocket)  
- Authentication: JWT tokens via Auth0  
- Static Content: Served from S3 behind CloudFront CDN  

## Directory Structure
```
├── cmd/                 # Entry points for services
├── internal/            # Application core logic
│   ├── api/             # HTTP handlers and routes
│   ├── service/         # Business logic and use cases
│   └── repository/      # Data persistence layer
├── web/                 # React application source
│   ├── components/      
│   ├── store/           
│   └── pages/           
├── scripts/             # CI/CD and deployment scripts
└── configs/             # Environment-specific configuration files
```

## Core Components
- **API Gateway**: Exposes RESTful endpoints for CRUD operations  
- **Task Service**: Manages task lifecycle and real-time updates  
- **User Service**: Handles authentication, authorization, and profile management  
- **Web Client**: SPA built with React for user interactions  

## Data Flow
1. Client makes HTTP/WebSocket calls to the API gateway  
2. API gateway routes requests to appropriate service  
3. Services interact with the database via repository layer  
4. Real-time updates push via WebSocket to connected clients  

## Infrastructure & Deployment
- CI/CD: GitHub Actions pipeline running lint, tests, and Docker image builds  
- Containerization: Docker for backend services  
- Orchestration: AWS ECS with Fargate  
- Secrets & Config: AWS Parameter Store and SSM  

## Conventions and Best Practices
- Follow 12-factor app methodology  
- Use environment variables for configuration  
- Write unit tests for service logic and integration tests for API routes  
- Apply consistent code formatting via Prettier (frontend) and go fmt (backend)  