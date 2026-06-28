package globale

var EnvPresets = map[string]string{
	"nodejs": `# App
NODE_ENV=development
PORT=3000
`,
	"go": `# App
GO_ENV=development
PORT=8080
`,
	"django": `# Django
DJANGO_SECRET_KEY=your-secret-key-here
DJANGO_DEBUG=True
DJANGO_ALLOWED_HOSTS=localhost,127.0.0.1
`,
	"database": `# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USER=postgres
DB_PASSWORD=password
DB_SSL_MODE=disable
`,
	"jwt": `# JWT Auth
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=7d
`,
	"smtp": `# SMTP Mail
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@email.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your@email.com
`,
	"s3": `# S3 / MinIO
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=my-bucket
S3_REGION=us-east-1
S3_USE_SSL=false
`,
	"docker": `# Docker
COMPOSE_PROJECT_NAME=myapp
DOCKER_REGISTRY=docker.io
`,
}
