package globale

var GitignoreTemplates = map[string]string{
	"Go": `# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.sum
vendor/
`,
	"Python": `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.env
.venv
env/
venv/
dist/
build/
*.egg-info/
.eggs/
`,
	"Node": `# Node
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.pnpm-debug.log*
dist/
build/
.cache/
`,
	"Java": `# Java
*.class
*.jar
*.war
*.ear
*.zip
*.tar.gz
target/
.mvn/
`,
	"Rust": `# Rust
/target/
Cargo.lock
`,
	"C++": `# C++
*.o
*.obj
*.exe
*.out
*.app
*.i*86
*.x86_64
*.hex
build/
`,
	"Angular": `# Angular
/dist/
/tmp/
/out-tsc/
.angular/
`,
	"React": `# React
/build/
/.next/
/out/
`,
	"NextJS": `# Next.js
/.next/
/out/
next-env.d.ts
`,
	"Vue": `# Vue
/dist/
`,
	"Laravel": `# Laravel
/vendor/
/node_modules/
/public/hot
/public/storage
/storage/*.key
.env
.env.backup
`,
	"Django": `# Django
*.log
local_settings.py
db.sqlite3
db.sqlite3-journal
media/
`,
	"Docker": `# Docker
.dockerignore
docker-compose.override.yml
`,
	"Terraform": `# Terraform
*.tfstate
*.tfstate.*
.terraform/
*.tfvars
crash.log
`,
	"Linux": `# Linux
*~
.fuse_hidden*
.directory
.Trash-*
.nfs*
`,
	"macOS": `# macOS
.DS_Store
.AppleDouble
.LSOverride
._*
.Spotlight-V100
.Trashes
`,
	"Windows": `# Windows
Thumbs.db
Thumbs.db:encryptable
ehthumbs.db
Desktop.ini
$RECYCLE.BIN/
`,
	"VSCode": `# VSCode
.vscode/*
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json
*.code-workspace
`,
	"JetBrains": `# JetBrains
.idea/
*.iml
*.iws
out/
`,
}

var CommonGitignore = `# Common
.env
.env.local
.env.*.local
*.log
logs/
tmp/
temp/
.cache/
coverage/
`
