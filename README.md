# MXShop Microservices Project

## Project Overview

MXShop is a modern e-commerce platform built using a microservices architecture with Go. The project is organized into separate services for different business domains, providing a scalable and maintainable solution for online retail operations.

## Project Structure

The project is divided into two main components:

### API Layer (`mxshop_api`)

Frontend-facing APIs that handle client requests and communicate with backend services.

- **goods_web**: API for product-related operations
- **order_web**: API for order management
- **oss_web**: API for object storage service
- And other web services...

### Service Layer (`mxshop_srv`)

Backend microservices that implement core business logic.

- **order_srv**: Order processing service
- And other backend services...

## Technologies

- **Go**: Primary programming language
- **gRPC**: For internal service communication
- **OpenTracing**: For distributed tracing (using otgrpc)
- **Object Storage**: For file uploads and management

## Getting Started

(Instructions for setting up development environment, running services, etc.)

## API Documentation

(Links or information about API documentation)

## Contributing

(Guidelines for contributing to the project)

## License

(License information)
