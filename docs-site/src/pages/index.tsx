import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import useBaseUrl from '@docusaurus/useBaseUrl';
import Layout from '@theme/Layout';
import Head from '@docusaurus/Head';
import styles from './index.module.css';

const DEMO_LINES = [
  { type: 'header',  text: '  oat-latte notes' },
  { type: 'divider', text: '' },
  { type: 'list',    active: true,  text: '  ▶  Meeting recap                                        2026-03-21' },
  { type: 'list',    active: false, text: '     Buy groceries                                         2026-03-20' },
  { type: 'list',    active: false, text: '     Refactor canvas focus logic                           2026-03-19' },
  { type: 'list',    active: false, text: '     Ship v0.1.0                                           2026-03-18' },
  { type: 'list',    active: false, text: '     Ideas for homepage demo                               2026-03-17' },
  { type: 'blank',   text: '' },
  { type: 'blank',   text: '' },
  { type: 'blank',   text: '' },
  { type: 'divider', text: '' },
  { type: 'footer',  text: '  n New    e Edit    Del Delete    Tab Next    q Quit' },
];

function TerminalDemo() {
  return (
    <div className={styles.terminal} aria-hidden="true">
      <div className={styles.terminalBar}>
        <span className={styles.dot}/>
        <span className={styles.dot}/>
        <span className={styles.dot}/>
        <span className={styles.terminalTitle}>oat-latte — notes</span>
      </div>
      <div className={styles.terminalBody}>
        {DEMO_LINES.map((line, i) => (
          <div
            key={i}
            className={[
              styles.termLine,
              line.type === 'header'  ? styles.lineHeader  : '',
              line.type === 'footer'  ? styles.lineFooter  : '',
              line.type === 'divider' ? styles.lineDivider : '',
              (line as any).active    ? styles.lineActive  : '',
            ].filter(Boolean).join(' ')}
          >
            {line.text}
          </div>
        ))}
      </div>
    </div>
  );
}

const FEATURES = [
  {
    icon: '⬡',
    title: 'Two-pass layout',
    body: 'Measure then Render. Every component declares its size before paint — no layout thrash, predictable pixel placement every frame.',
  },
  {
    icon: '◈',
    title: 'Composable widgets',
    body: 'Text, Button, EditText, List, CheckBox, ProgressBar, NotificationManager — all themed, all focusable, all composable via VBox/HBox/Grid.',
  },
  {
    icon: '◉',
    title: 'Focus-first design',
    body: 'DFS focus collection, Tab/Shift-Tab cycling, proxy pattern for key interception, and FocusByRef for programmatic jumps.',
  },
  {
    icon: '◐',
    title: 'Five built-in themes',
    body: 'Default (ANSI-16), Dark, Light, Dracula, Nord. Apply once with WithTheme — every widget inherits and can override via Style.Merge.',
  },
  {
    icon: '◇',
    title: 'Modal dialogs',
    body: 'ShowDialog stacks overlays with a full-screen scrim. Dialogs size by fixed cells or percentage of terminal. HideDialog restores focus.',
  },
  {
    icon: '◑',
    title: 'Goroutine-safe redraws',
    body: 'NotifyChannel lets background goroutines trigger repaints without locks. Key handlers run on the main goroutine — no races.',
  },
  {
    icon: '◭',
    title: 'Automatic focus highlight',
    body: 'Border panels light up with the theme\'s focus colour the moment any descendant receives focus. Zero configuration — it just works.',
  },
  {
    icon: '◮',
    title: 'Context-aware tab gating',
    body: 'FocusGuard lets any component exclude its entire subtree from Tab cycling at runtime. Build mode-gated panels without managing focus lists manually.',
  },
];

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  const version = (siteConfig.customFields?.version as string) ?? 'v0.1.0';
  const logoUrl = useBaseUrl('/img/logo.png');
  const ogImage = `${siteConfig.url}/img/android-chrome-512x512.png`;
  const canonicalUrl = siteConfig.url;
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <Head>
        {/* Primary */}
        <meta name="description" content={siteConfig.tagline} />
        <link rel="canonical" href={canonicalUrl} />

        {/* OpenGraph */}
        <meta property="og:type" content="website" />
        <meta property="og:url" content={canonicalUrl} />
        <meta property="og:site_name" content={siteConfig.title} />
        <meta property="og:title" content={`${siteConfig.title} — ${siteConfig.tagline}`} />
        <meta property="og:description" content="A component-based TUI framework for Go. Two-pass layout, composable widgets, five built-in themes, modal dialogs, and goroutine-safe redraws." />
        <meta property="og:image" content={ogImage} />
        <meta property="og:image:type" content="image/png" />
        <meta property="og:image:width" content="512" />
        <meta property="og:image:height" content="512" />
        <meta property="og:image:alt" content="oat-latte logo" />

        {/* Twitter / X Card */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content={`${siteConfig.title} — ${siteConfig.tagline}`} />
        <meta name="twitter:description" content="A component-based TUI framework for Go. Two-pass layout, composable widgets, five built-in themes, modal dialogs, and goroutine-safe redraws." />
        <meta name="twitter:image" content={ogImage} />
        <meta name="twitter:image:alt" content="oat-latte logo" />
      </Head>
      <main className={styles.page}>

        {/* hero */}
        <section className={styles.hero}>
          <div className={styles.heroText}>
            <div className={styles.badge}>Go · TUI · component-based</div>
            <div className={styles.heroTitleRow}>
              <img src={logoUrl} alt="oat-latte" className={styles.heroLogo} />
              <h1 className={styles.heroTitle}>
                <span className={styles.heroWord}>oat</span>
                <span className={styles.herySep}>-</span>
                <span className={styles.heroWord}>latte</span>
              </h1>
            </div>
            <p className={styles.heroSub}>
              A component-based TUI framework for Go.<br/>
              Measure. Render. Ship.
            </p>
            <div className={styles.heroCta}>
              <Link className={styles.ctaPrimary} to="/docs/intro">
                Get started
              </Link>
              <Link className={styles.ctaSecondary} to="/docs/concepts">
                Core concepts
              </Link>
              <a
                className={styles.ctaGhost}
                href="https://github.com/antoniocali/oat-latte"
                target="_blank"
                rel="noopener noreferrer"
              >
                GitHub ↗
              </a>
            </div>
            <div className={styles.installBlock}>
              <span className={styles.installPrompt}>$</span>
              <span className={styles.installCmd}>go get github.com/antoniocali/oat-latte@{version}</span>
            </div>
          </div>
          <div className={styles.heroDemo}>
            <TerminalDemo />
          </div>
        </section>

        {/* feature grid */}
        <section className={styles.features}>
          <div className={styles.featuresGrid}>
            {FEATURES.map((f) => (
              <div key={f.title} className={styles.featureCard}>
                <span className={styles.featureIcon}>{f.icon}</span>
                <h3 className={styles.featureTitle}>{f.title}</h3>
                <p className={styles.featureBody}>{f.body}</p>
              </div>
            ))}
          </div>
        </section>

        {/* code sample */}
        <section className={styles.codeSample}>
          <h2 className={styles.codeSampleHeading}>From zero to running in minutes</h2>
          <pre className={styles.codeBlock}>{`package main

import (
    "log"
    oat "github.com/antoniocali/oat-latte"
    "github.com/antoniocali/oat-latte/latte"
    "github.com/antoniocali/oat-latte/layout"
    "github.com/antoniocali/oat-latte/widget"
)

func main() {
    input := widget.NewEditText().
        WithHint("Name").
        WithPlaceholder("Type something…")

    btn := widget.NewButton("Say hello", func() {
        // handle press
    })

    body := layout.NewBorder(
        layout.NewVBox(input, btn),
    ).WithTitle("Hello")

    app := oat.NewCanvas(
        oat.WithTheme(latte.ThemeDark),
        oat.WithBody(body),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}`}</pre>
        </section>

      </main>
    </Layout>
  );
}
