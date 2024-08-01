# simple-student-docker-project

This project is a comprehensive demonstration of an application that integrates RabbitMQ, PostgreSQL, and Go. 
The system includes a reader that reads data from a text file and publish to RabbitMQ, a writer that comsume this data from RabbitMQ to a PostgreSQL database, and an analyser that retrieves and processes this data to another file.

## Architecture

The architecture consists of several Dockerized services:
- **RabbitMQ**: Message broker to handle data communication between services.
- **PostgreSQL**: Stores students data.
- **Reader**: Reads data from a text file and publish to RabbitMQ.
- **Writer**: Consumes data from RabbitMQ and writes it to the PostgreSQL database.
- **Analyser**: Retrieves and processes students data from PostgreSQL into another file.

## Prerequisites

- Docker
- Docker Compose
