package globale

import "fmt"

var ComposeServices = map[string]func(string) string{
	"postgres": func(version string) string {
		v := version
		if v == "" {
			v = "16"
		}
		return fmt.Sprintf(`  postgres:
    image: postgres:%s
    container_name: postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-password}
      POSTGRES_DB: ${DB_NAME:-mydb}
    ports:
      - "127.0.0.1:5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
`, v)
	},

	"mysql": func(version string) string {
		v := version
		if v == "" {
			v = "8"
		}
		return fmt.Sprintf(`  mysql:
    image: mysql:%s
    container_name: mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD:-rootpassword}
      MYSQL_DATABASE: ${DB_NAME:-mydb}
      MYSQL_USER: ${DB_USER:-user}
      MYSQL_PASSWORD: ${DB_PASSWORD:-password}
    ports:
      - "127.0.0.1:3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
`, v)
	},

	"redis": func(version string) string {
		v := version
		if v == "" {
			v = "7"
		}
		return fmt.Sprintf(`  redis:
    image: redis:%s-alpine
    container_name: redis
    restart: unless-stopped
    ports:
      - "127.0.0.1:6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
`, v)
	},

	"mongodb": func(version string) string {
		v := version
		if v == "" {
			v = "7"
		}
		return fmt.Sprintf(`  mongodb:
    image: mongo:%s
    container_name: mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER:-admin}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD:-password}
    ports:
      - "127.0.0.1:27017:27017"
    volumes:
      - mongodb_data:/data/db
`, v)
	},

	"nginx": func(_ string) string {
		return `  nginx:
    image: nginx:alpine
    container_name: nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
`
	},

	"caddy": func(_ string) string {
		return `  caddy:
    image: caddy:latest
    container_name: caddy
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
`
	},

	"minio": func(_ string) string {
		return `  minio:
    image: minio/minio:latest
    container_name: minio
    restart: unless-stopped
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: ${MINIO_USER:-minioadmin}
      MINIO_ROOT_PASSWORD: ${MINIO_PASSWORD:-minioadmin}
    ports:
      - "127.0.0.1:9000:9000"
      - "127.0.0.1:9001:9001"
    volumes:
      - minio_data:/data
`
	},

	"rabbitmq": func(_ string) string {
		return `  rabbitmq:
    image: rabbitmq:management-alpine
    container_name: rabbitmq
    restart: unless-stopped
    environment:
      RABBITMQ_DEFAULT_USER: ${RABBITMQ_USER:-admin}
      RABBITMQ_DEFAULT_PASS: ${RABBITMQ_PASSWORD:-password}
    ports:
      - "127.0.0.1:5672:5672"
      - "127.0.0.1:15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
`
	},

	"grafana": func(_ string) string {
		return `  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    restart: unless-stopped
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_PASSWORD:-admin}
    ports:
      - "127.0.0.1:3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
`
	},

	"prometheus": func(_ string) string {
		return `  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    restart: unless-stopped
    ports:
      - "127.0.0.1:9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
`
	},

	"portainer": func(_ string) string {
		return `  portainer:
    image: portainer/portainer-ce:lts
    container_name: portainer
    restart: unless-stopped
    ports:
      - "127.0.0.1:9000:9000"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - portainer_data:/data
`
	},
}

var VolumeNames = map[string]string{
	"postgres":   "postgres_data",
	"mysql":      "mysql_data",
	"redis":      "redis_data",
	"mongodb":    "mongodb_data",
	"minio":      "minio_data",
	"rabbitmq":   "rabbitmq_data",
	"grafana":    "grafana_data",
	"prometheus": "prometheus_data",
	"portainer":  "portainer_data",
	"caddy":      "caddy_data\n  caddy_config",
}
