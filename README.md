# LLM-Golang

This project is a web application for interacting with a large language model (LLM). It features user authentication with GitHub OAuth, a responsive UI built with Bootstrap and HTMX, and integration with a remote PostgreSQL database hosted on Neon Console.

---

## Features

- **User Authentication**: Login via GitHub OAuth.
- **Prompt Submission**: Submit prompts to an LLM and view responses.
- **Profile Management**: User profiles with avatar, email, and username display.
- **Database Integration**: Save user information, prompts, and responses to a PostgreSQL database.
- **Dockerized Setup**: Fully containerized for consistent deployment.
- **Secure Configuration**: Environment variables and Docker Secrets for sensitive data.

---

## Prerequisites

- Docker and Docker Compose installed.
- Access to a PostgreSQL database (e.g., Neon Console).
- A GitHub OAuth App for authentication:
  - [Create a GitHub OAuth App](https://docs.github.com/en/developers/apps/building-oauth-apps/creating-an-oauth-app)

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-repo/LLM-Golang.git
cd LLM-Golang
