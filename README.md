# Great Nigeria Library Foundation

**Open-source platform for educational and cultural content management**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-blue.svg)](https://reactjs.org)

## 🎯 What is the Foundation?

The Great Nigeria Library Foundation is a complete, production-ready platform for building educational and cultural content platforms. It provides all the essential features you need to create your own library, educational platform, or community content hub.

## ✨ Features

### 🔐 Authentication & User Management
- Secure user registration and login
- JWT-based authentication
- User profiles and settings
- Password reset functionality

### 📚 Content Management
- Book and document reading
- Demo content included
- Search and discovery
- Content organization

### 💬 Community Features
- Discussion forums
- User interactions
- Community moderation
- Social features

### 🛠️ Developer-Friendly
- RESTful API design
- Comprehensive documentation
- Docker deployment
- Extensible architecture

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/yerenwgventures/GreatNigeriaLibrary.git
cd GreatNigeriaLibrary

# Start the foundation platform
docker-compose up -d

# Access the platform
open http://localhost:8080
```

### First Steps

1. **Register an account** at http://localhost:8080/register
2. **Explore demo content** in the books section
3. **Join discussions** in the community forums
4. **Read the platform guide** for detailed instructions

## 📖 Demo Content

The foundation includes sample content to help you get started:

- **Platform User Guide** - Complete guide to using the platform
- **Nigerian History Overview** - Educational content example
- **Community Guidelines** - Best practices for community engagement

## 🏗️ Architecture

### Backend (Go)
```
foundation/backend/
├── cmd/                    # Service entry points
├── internal/               # Core business logic
│   ├── auth/              # Authentication
│   ├── content/           # Content management
│   ├── discussion/        # Forums
│   └── groups/            # Community groups
├── pkg/                   # Shared packages
│   ├── models/            # Data models
│   └── common/            # Utilities
└── main.go               # Application entry point
```

### Frontend (React + TypeScript)
```
foundation/frontend/
├── src/
│   ├── features/          # Feature modules
│   │   ├── auth/         # Authentication UI
│   │   ├── books/        # Content reading
│   │   ├── forum/        # Discussion UI
│   │   └── search/       # Search interface
│   ├── components/       # Reusable components
│   ├── types/           # TypeScript definitions
│   └── api/             # API integration
```

## 🔧 Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=foundation_user
DB_PASSWORD=foundation_pass
DB_DATABASE=great_nigeria_foundation

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Authentication
JWT_SECRET=your-secret-key
ACCESS_TOKEN_EXPIRATION=15m
REFRESH_TOKEN_EXPIRATION=168h

# Server
SERVER_PORT=8080
ENVIRONMENT=development
```

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update profile

### Content Endpoints
- `GET /api/v1/content/books` - List available books
- `GET /api/v1/content/books/:id` - Get specific book
- `GET /api/v1/content/search` - Search content

### Discussion Endpoints
- `GET /api/v1/discussion/forums` - List forums
- `POST /api/v1/discussion/topics` - Create topic
- `GET /api/v1/discussion/topics/:id` - Get topic

## 🧪 Testing

```bash
# Run backend tests
cd foundation/backend
go test ./...

# Run frontend tests
cd foundation/frontend
npm test
```

## 🚀 Deployment

### Docker Deployment
```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Deployment
```bash
# Build backend
cd foundation/backend
go build -o app ./main.go

# Build frontend
cd foundation/frontend
npm run build

# Deploy to your server
```

## 🔌 Extending the Platform

The foundation is designed to be extended with additional features:

### Premium Features (Available Separately)
- Payment processing and e-commerce
- Live streaming and events
- Advanced analytics and reporting
- AI-powered recommendations
- Gamification and achievements

### Custom Integrations
- Third-party authentication (OAuth)
- External content APIs
- Custom themes and branding
- Mobile app development

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup
```bash
# Clone and setup
git clone https://github.com/yerenwgventures/GreatNigeriaLibrary.git
cd GreatNigeriaLibrary

# Install dependencies
cd foundation/backend && go mod download
cd ../frontend && npm install

# Start development servers
docker-compose -f docker-compose.dev.yml up
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Acknowledgments

- Built with ❤️ for the Nigerian educational community
- Inspired by the need for accessible, quality educational platforms
- Designed to empower developers and educators across Africa

## 📞 Support

- **Documentation**: [docs.greatnigeria.com](https://docs.greatnigeria.com)
- **Community Forum**: [community.greatnigeria.com](https://community.greatnigeria.com)
- **Issues**: [GitHub Issues](https://github.com/yerenwgventures/GreatNigeriaLibrary/issues)
- **Email**: support@greatnigeria.com

---

**Start building your educational platform today with the Great Nigeria Library Foundation!**
