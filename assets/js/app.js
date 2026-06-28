// ─── Utilitaires ────────────────────────────────────────────────────────────

const langMap = {
    gitignore:   'plaintext',
    license:     'plaintext',
    packagejson: 'json',
    env:         'typescript',
    dockerfile:  'typescript',
}

function showOutput(content) {
    const output = document.getElementById('output')
    const outputContent = document.getElementById('outputContent')
    if (!output || !outputContent) return
    output.style.display = 'block'

    const page = document.body.dataset.page
    const lang = langMap[page] || 'plaintext'

    outputContent.removeAttribute('data-highlighted')
    outputContent.className = 'language-' + lang
    outputContent.textContent = content

    if (window.hljs) {
        hljs.highlightElement(outputContent)
    }
}

function setupCopy() {
    const copyBtn = document.getElementById('copyBtn')
    const outputContent = document.getElementById('outputContent')
    if (!copyBtn || !outputContent) return

    copyBtn.addEventListener('click', () => {
        navigator.clipboard.writeText(outputContent.textContent).then(() => {
            copyBtn.textContent = 'Copied!'
            copyBtn.classList.add('copied')
            setTimeout(() => {
                copyBtn.textContent = 'Copy'
                copyBtn.classList.remove('copied')
            }, 2000)
        })
    })
}

async function postAndShow(url, body) {
    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        })

        if (!res.ok) {
            console.error('HTTP error:', res.status)
            return
        }

        const data = await res.json()
        if (data.content) {
            showOutput(data.content)
        } else if (data.error) {
            alert('Error: ' + data.error)
        }
    } catch (err) {
        console.error('postAndShow error:', err)
    }
}

// ─── Tags multi-select ───────────────────────────────────────────────────────

function setupTags(generateFn) {
    const tags = document.querySelectorAll('.tag')
    const generateBtn = document.getElementById('generateBtn')
    const selectedTagsContainer = document.getElementById('selectedTags')
    const selected = new Set()

    tags.forEach(tag => {
        tag.addEventListener('click', () => {
            const value = tag.dataset.value
            if (!value) return

            if (selected.has(value)) {
                selected.delete(value)
                tag.classList.remove('selected')
            } else {
                selected.add(value)
                tag.classList.add('selected')
            }

            if (selectedTagsContainer) {
                if (selected.size === 0) {
                    selectedTagsContainer.innerHTML = '<p class="empty-state">No technologies selected yet</p>'
                } else {
                    selectedTagsContainer.innerHTML = [...selected]
                        .map(v => `<span class="selected-tag">${v}</span>`)
                        .join('')
                }
            }

            if (generateBtn) {
                generateBtn.disabled = selected.size === 0
            }
        })
    })

    if (generateBtn && generateFn) {
        generateBtn.addEventListener('click', () => {
            const values = [...selected]
            if (values.length === 0) return
            generateFn(values)
        })
    }
}

// ─── Init par page ───────────────────────────────────────────────────────────

const page = document.body.dataset.page

if (page === 'gitignore') {
    setupTags(async (technologies) => {
        await postAndShow('/api/gitignore', { technologies })
    })
    setupCopy()
}

if (page === 'license') {
    const btn = document.getElementById('generateBtn')
    if (btn) {
        btn.addEventListener('click', async () => {
            await postAndShow('/api/license', {
                type:    document.getElementById('licenseType')?.value || 'mit',
                author:  document.getElementById('authorName')?.value || '',
                year:    document.getElementById('year')?.value || '2026',
                project: document.getElementById('projectName')?.value || '',
            })
        })
    }
    setupCopy()
}

if (page === 'packagejson') {
    const btn = document.getElementById('generateBtn')
    if (btn) {
        btn.addEventListener('click', async () => {
            const pkgType = document.querySelector('input[name="pkgType"]:checked')
            await postAndShow('/api/packagejson', {
                name:           document.getElementById('pkgName')?.value || '',
                version:        document.getElementById('pkgVersion')?.value || '1.0.0',
                description:    document.getElementById('pkgDesc')?.value || '',
                author:         document.getElementById('pkgAuthor')?.value || '',
                license:        document.getElementById('pkgLicense')?.value || 'MIT',
                type:           pkgType?.value || '',
                includeScripts: document.getElementById('includeScripts')?.checked ?? false,
                includeEngines: document.getElementById('includeEngines')?.checked ?? false,
                includePrivate: document.getElementById('includePrivate')?.checked ?? false,
            })
        })
    }
    setupCopy()
}

if (page === 'env') {
    setupTags(async (presets) => {
        await postAndShow('/api/env', {
            presets,
            appName: document.getElementById('appName')?.value || '',
        })
    })
    setupCopy()
}

if (page === 'dockerfile') {
    const btn = document.getElementById('generateBtn')
    if (btn) {
        btn.addEventListener('click', async () => {
            await postAndShow('/api/dockerfile', {
                lang:        document.getElementById('lang')?.value || 'go',
                langVersion: document.getElementById('langVersion')?.value || '',
                port:        document.getElementById('port')?.value || '8080',
                workdir:     document.getElementById('workdir')?.value || '/app',
                multistage:  document.getElementById('multistage')?.checked ?? true,
                nonroot:     document.getElementById('nonroot')?.checked ?? true,
                healthcheck: document.getElementById('healthcheck')?.checked ?? false,
            })
        })
    }
    setupCopy()
}