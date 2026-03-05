# inge.4pr.es — blog

Source for [inge.4pr.es](https://inge.4pr.es), a personal blog about software engineering, observability, and performance. Built with [Hugo](https://gohugo.io) and a custom theme called **signal**.

## Requirements

- [Hugo extended](https://gohugo.io/installation/) v0.120+ (extended edition required for Dart Sass)
- [Dart Sass](https://sass-lang.com/install/) (`sass` binary in PATH)

Check your Hugo edition:
```bash
hugo env | head -1
# should include "extended" or list libwebp/libvips
```

## Local development

```bash
cd inge.4pr.es
hugo server -D
```

Opens at http://localhost:1313. The `-D` flag renders draft posts.

To do a full production build locally (output goes to `public/`, not tracked in git):
```bash
cd inge.4pr.es
hugo -e production
```

## Writing a post

Create a Markdown file in `inge.4pr.es/content/post/`:

```yaml
---
title: "Post Title"
author: "inge4pres"
date: 2025-01-01T12:00:00+01:00
subtitle: "A short description shown on the post card"
categories:
- tech
tags:
- golang
- observability
draft: true
---

Content here...
```

Set `draft: false` (or remove the field) when ready to publish. The permalink is derived from the filename, so `my-post.md` becomes `inge.4pr.es/my-post`.

## Theme: signal

The `signal` theme lives entirely in `inge.4pr.es/themes/signal/`. It is a first-party theme — no git submodules.

### Design principles

- **CSS custom properties** drive all colours, spacing, and typography. Tokens are defined in `assets/scss/_tokens.scss`.
- **Dark/light mode** is based on a `data-theme` attribute on `<html>`. On page load, an inline script in `<head>` reads `localStorage` (falling back to `prefers-color-scheme`) and sets the attribute before the first paint — no flash. A toggle button in the navbar lets the user override and persists the choice.
- **System fonts only** — no external font requests, zero latency penalty.
- **SCSS compiled by Hugo Pipes** via `resources.ToCSS` with Dart Sass. No Node.js or npm required.

### File layout

```
themes/signal/
├── assets/
│   └── scss/
│       ├── main.scss               # entry point — imports all partials
│       ├── _tokens.scss            # design tokens (light default + [data-theme="dark"])
│       ├── _reset.scss             # modern CSS reset
│       ├── _typography.scss        # heading scale, .prose wrapper styles
│       ├── _layout.scss            # containers, post grid, hero, breakpoints
│       ├── _syntax.scss            # code block wrapper + language label
│       ├── _syntax-chroma.scss     # generated Chroma CSS (dracula palette)
│       └── _components/
│           ├── _header.scss        # sticky nav, dropdown, theme toggle
│           ├── _footer.scss        # three-column footer, SVG social icons
│           ├── _post-card.scss     # card used in grids
│           ├── _post-single.scss   # post page, meta, social share, comments
│           ├── _tags.scss          # tag pills and terms-list
│           └── _pagination.scss
└── layouts/
    ├── _default/
    │   ├── baseof.html             # base shell
    │   ├── single.html             # post and static pages
    │   ├── list.html               # section and taxonomy pages
    │   └── terms.html              # tags index
    ├── index.html                  # homepage (hero + post grid)
    ├── 404.html
    └── partials/
        ├── head.html               # <head>, CSS pipeline, OG tags, theme-init script
        ├── header.html             # navbar with avatar, nav links, theme toggle
        ├── footer.html             # copyright, RSS, social icons
        ├── post-card.html          # card component
        ├── post-meta.html          # date · author · reading time · word count
        ├── related-posts.html      # up to 3 related posts
        ├── social-share.html       # Twitter, LinkedIn, copy-link
        ├── comments.html           # Disqus (guarded by DisqusShortname config)
        └── analytics.html          # Google Analytics gtag (skipped in dev)
```

### Regenerating syntax highlight CSS

If you change the Chroma style, regenerate the CSS from `inge.4pr.es/`:

```bash
hugo gen chromastyles --style=dracula > themes/signal/assets/scss/_syntax-chroma.scss
```

### Key Hugo config (`beautifulhugo.config.toml`)

| Setting | Value | Reason |
|---|---|---|
| `theme` | `signal` | active theme |
| `[markup.goldmark.renderer] unsafe` | `true` | older posts contain raw HTML |
| `[markup.highlight] noClasses` | `false` | class-based Chroma output (CSS-driven) |
| `[markup.highlight] style` | `dracula` | syntax palette |
| `[permalinks] post` | `/:slugorcontentbasename` | clean URLs from filenames |

## Deployment

The site is deployed through **Cloudflare Pages**. Pushing to the `master` branch triggers a build automatically.
