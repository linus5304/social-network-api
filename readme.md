# GopherSocial API

A robust social networking REST API built with Go, featuring user authentication, post management, and real-time interactions.

## Overview

GopherSocial is a feature-rich social networking API that enables users to create posts, follow other users, and engage through comments. Built with modern Go practices and following clean architecture principles.

### API Documentation

Explore and test the API endpoints using our interactive Swagger documentation:
[Live API Documentation](https://social-network-api-634079758108.us-central1.run.app/v1/swagger/index.html)

### Key Features

- User authentication with JWT
- Post creation and management
- User following system
- Comment system
- Real-time feed updates
- Full text search capabilities
- Rate limiting
- Redis caching
- Email notifications
- Swagger/OpenAPI documentation

### Tech Stack

- **Language:** Go 1.23
- **Database:** PostgreSQL with full-text search
- **Caching:** Redis
- **Documentation:** Swagger/OpenAPI
- **Authentication:** JWT
- **Logging:** Uber's zap logger
- **Testing:** Go testing package with testify
- **CI/CD:** GitHub Actions

## Getting Started

Prerequisites:

- Go 1.23+
- PostgreSQL
- Redis (optional)
- Make
