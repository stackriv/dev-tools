// ─── Utilitaires ────────────────────────────────────────────────────────────

const langMap = {
    gitignore:   'plaintext',
    license:     'plaintext',
    packagejson: 'json',
    env:         'typescript',
    dockerfile:  'typescript',
};

function showOutput(content) {
    const output = document.getElementById('output');
    const outputContent = document.getElementById('outputContent');
    if (!output || !outputContent) return;
    output.style.display = 'block';

    const page = document.body.dataset.page;
    const lang = langMap[page] || 'plaintext';

    outputContent.removeAttribute('data-highlighted');
    outputContent.className = 'language-' + lang;
    outputContent.textContent = content;

    if (window.hljs) {
        hljs.highlightElement(outputContent);
    }
}

function setupCopy() {
    const copyBtn = document.getElementById('copyBtn');
    const outputContent = document.getElementById('outputContent');
    if (!copyBtn || !outputContent) return;

    copyBtn.addEventListener('click', () => {
        navigator.clipboard.writeText(outputContent.textContent).then(() => {
            copyBtn.textContent = 'Copied!';
            copyBtn.classList.add('copied');
            setTimeout(() => {
                copyBtn.textContent = 'Copy';
                copyBtn.classList.remove('copied');
            }, 2000);
        });
    });
}

async function postAndShow(url, body) {
    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
        });

        if (!res.ok) {
            console.error('HTTP error:', res.status);
            return;
        }

        const data = await res.json();
        if (data.content) {
            showOutput(data.content);
        } else if (data.error) {
            alert('Error: ' + data.error);
        }
    } catch (err) {
        console.error('postAndShow error:', err);
    }
}

// ─── Tags multi-select ───────────────────────────────────────────────────────

function setupTags(generateFn) {
    const tags = document.querySelectorAll('.tag');
    const generateBtn = document.getElementById('generateBtn');
    const selectedTagsContainer = document.getElementById('selectedTags');
    const selected = new Set();

    tags.forEach(tag => {
        tag.addEventListener('click', () => {
            const value = tag.dataset.value;
            if (!value) return;

            if (selected.has(value)) {
                selected.delete(value);
                tag.classList.remove('selected');
            } else {
                selected.add(value);
                tag.classList.add('selected');
            }

            if (selectedTagsContainer) {
                if (selected.size === 0) {
                    selectedTagsContainer.innerHTML = '<p class="empty-state">No technologies selected yet</p>';
                } else {
                    selectedTagsContainer.innerHTML = [...selected]
                        .map(v => `<span class="selected-tag">${v}</span>`)
                        .join('');
                }
            }

            if (generateBtn) {
                generateBtn.disabled = selected.size === 0;
            }
        });
    });

    if (generateBtn && generateFn) {
        generateBtn.addEventListener('click', () => {
            const values = [...selected];
            if (values.length === 0) return;
            generateFn(values);
        });
    }
}

// ─── Init par page ───────────────────────────────────────────────────────────

const page = document.body.dataset.page;

if (page === 'gitignore') {
    setupTags(async (technologies) => {
        await postAndShow('/api/gitignore', { technologies });
    });
    setupCopy();
}

if (page === 'license') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            await postAndShow('/api/license', {
                type:    document.getElementById('licenseType')?.value || 'mit',
                author:  document.getElementById('authorName')?.value || '',
                year:    document.getElementById('year')?.value || '2026',
                project: document.getElementById('projectName')?.value || '',
            });
        });
    }
    setupCopy();
}

if (page === 'packagejson') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const pkgType = document.querySelector('input[name="pkgType"]:checked');
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
            });
        });
    }
    setupCopy();
}

if (page === 'env') {
    setupTags(async (presets) => {
        await postAndShow('/api/env', {
            presets,
            appName: document.getElementById('appName')?.value || '',
        });
    });
    setupCopy();
}

if (page === 'dockerfile') {
    const btn = document.getElementById('generateBtn');
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
            });
        });
    }
    setupCopy();
}

// ─── UUID ────────────────────────────────────────────────────────────────────

if (page === 'uuid') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const count = document.getElementById('uuidCount')?.value || 1;
            const format = document.querySelector('input[name="uuidFormat"]:checked')?.value || 'standard';

            const res = await fetch(`/api/uuid?count=${count}`);
            const data = await res.json();

            if (!data.uuids) return;

            let uuids = data.uuids;
            if (format === 'upper') uuids = uuids.map(u => u.toUpperCase());
            if (format === 'nodash') uuids = uuids.map(u => u.replace(/-/g, ''));

            showOutput(uuids.join('\n'));
        });
    }
    setupCopy();
}

// ─── Base64 ──────────────────────────────────────────────────────────────────

if (page === 'base64') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const input = document.getElementById('b64Input')?.value || '';
            const action = document.querySelector('input[name="b64Action"]:checked')?.value || 'encode';

            if (!input) return;

            try {
                const res = await fetch('/api/base64', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ input, action }),
                });
                const data = await res.json();
                if (data.result !== undefined) {
                    showOutput(data.result);
                } else if (data.error) {
                    alert('Error: ' + data.error);
                }
            } catch (err) {
                console.error(err);
            }
        });
    }
    setupCopy();
}

// ─── Hash ─────────────────────────────────────────────────────────────────────

if (page === 'hash') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const input = document.getElementById('hashInput')?.value || '';
            if (!input) return;

            try {
                const res = await fetch('/api/hash', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ input }),
                });
                const data = await res.json();

                const output = document.getElementById('output');
                if (output) output.style.display = 'block';

                const setValue = (id, val) => {
                    const el = document.getElementById(id);
                    if (el) el.textContent = val;
                };

                setValue('hashMd5', data.md5 || '');
                setValue('hashSha1', data.sha1 || '');
                setValue('hashSha256', data.sha256 || '');
                setValue('hashSha512', data.sha512 || '');

                // Copy buttons pour chaque hash
                document.querySelectorAll('.btn-copy-hash').forEach(copyBtn => {
                    copyBtn.onclick = () => {
                        const targetId = copyBtn.dataset.target;
                        const text = document.getElementById(targetId)?.textContent || '';
                        navigator.clipboard.writeText(text).then(() => {
                            copyBtn.textContent = 'Copied!';
                            copyBtn.classList.add('copied');
                            setTimeout(() => {
                                copyBtn.textContent = 'Copy';
                                copyBtn.classList.remove('copied');
                            }, 2000);
                        });
                    };
                });
            } catch (err) {
                console.error(err);
            }
        });
    }
}

// ─── JWT ─────────────────────────────────────────────────────────────────────

if (page === 'jwt') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const token = document.getElementById('jwtInput')?.value || '';
            if (!token) return;

            try {
                const res = await fetch('/api/jwt', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ token }),
                });
                const data = await res.json();

                if (data.error) {
                    alert('Error: ' + data.error);
                    return;
                }

                const output = document.getElementById('output');
                if (output) output.style.display = 'block';

                const setAndHighlight = (id, content) => {
                    const el = document.getElementById(id);
                    if (!el) return;
                    el.textContent = content;
                    if (window.hljs && el.classList.contains('language-json')) {
                        el.removeAttribute('data-highlighted');
                        hljs.highlightElement(el);
                    }
                };

                setAndHighlight('jwtHeader', data.header || '');
                setAndHighlight('jwtPayload', data.payload || '');

                const sig = document.getElementById('jwtSignature');
                if (sig) sig.textContent = data.signature || '';

                // Copy buttons
                document.querySelectorAll('.btn-copy-hash').forEach(copyBtn => {
                    copyBtn.onclick = () => {
                        const targetId = copyBtn.dataset.target;
                        const text = document.getElementById(targetId)?.textContent || '';
                        navigator.clipboard.writeText(text).then(() => {
                            copyBtn.textContent = 'Copied!';
                            copyBtn.classList.add('copied');
                            setTimeout(() => {
                                copyBtn.textContent = 'Copy';
                                copyBtn.classList.remove('copied');
                            }, 2000);
                        });
                    };
                });
            } catch (err) {
                console.error(err);
            }
        });
    }
}

// ─── Regex ───────────────────────────────────────────────────────────────────

if (page === 'regex') {
    const btn = document.getElementById('generateBtn');
    if (btn) {
        btn.addEventListener('click', async () => {
            const pattern = document.getElementById('regexPattern')?.value || '';
            const input   = document.getElementById('regexInput')?.value || '';
            if (!pattern) return;

            try {
                const res = await fetch('/api/regex', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        pattern,
                        input,
                        flags: {
                            global:      document.getElementById('flagGlobal')?.checked ?? true,
                            insensitive: document.getElementById('flagInsensitive')?.checked ?? false,
                            multiline:   document.getElementById('flagMultiline')?.checked ?? false,
                        },
                    }),
                });
                const data = await res.json();

                if (data.error) {
                    alert('Error: ' + data.error);
                    return;
                }

                const output = document.getElementById('output');
                if (output) output.style.display = 'block';

                const countEl = document.getElementById('regexCount');
                if (countEl) {
                    countEl.textContent = `${data.count} match${data.count !== 1 ? 'es' : ''}`;
                }

                const matchesEl = document.getElementById('regexMatches');
                if (!matchesEl) return;

                if (!data.matches || data.matches.length === 0) {
                    matchesEl.innerHTML = '<p style="color:#64748b;font-size:0.9rem;">No matches found.</p>';
                    return;
                }

                matchesEl.innerHTML = data.matches.map((m, i) => `
                    <div class="regex-match">
                        <div class="regex-match-header">
                            <span>Match ${i + 1}</span>
                            <span>position ${m.start}–${m.end}</span>
                        </div>
                        <span class="regex-match-value">${escapeHtml(m.match)}</span>
                        ${m.groups && m.groups.length > 0 ? `
                            <div class="regex-groups">
                                Groups: ${m.groups.map((g, gi) => `<span class="regex-group">$${gi + 1}: "${escapeHtml(g)}"</span>`).join(', ')}
                            </div>
                        ` : ''}
                    </div>
                `).join('');
            } catch (err) {
                console.error(err);
            }
        });
    }
}

function escapeHtml(str) {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}
