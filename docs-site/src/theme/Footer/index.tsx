/**
 * Swizzled Footer — oat-latte
 *
 * Design language: full B&W, 2px solid borders, IBM Plex Mono, invert-on-hover.
 * Replaces the Docusaurus default footer entirely.
 */

import React from 'react';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './styles.module.css';

export default function Footer(): React.ReactElement | null {
  const { siteConfig } = useDocusaurusContext();
  const year = new Date().getFullYear();

  return (
    <footer className={styles.footer}>
      <div className={styles.inner}>

        {/* Left — brand block */}
        <div className={styles.brand}>
          <span className={styles.brandName}>{siteConfig.title}</span>
          <p className={styles.tagline}>{siteConfig.tagline}</p>
          <p className={styles.copy}>
            &copy; {year} oat-latte. MIT licence.
          </p>
        </div>

        {/* Middle — docs links */}
        <nav className={styles.col} aria-label="Documentation links">
          <span className={styles.colHeading}>Docs</span>
          <ul className={styles.linkList}>
            <li><Link className={styles.link} to="/docs/intro">Getting Started</Link></li>
            <li><Link className={styles.link} to="/docs/concepts">Core Concepts</Link></li>
            <li><Link className={styles.link} to="/docs/widgets/text">Widgets</Link></li>
            <li><Link className={styles.link} to="/docs/layout">Layout</Link></li>
          </ul>
        </nav>

        {/* Right — external links */}
        <nav className={styles.col} aria-label="External links">
          <span className={styles.colHeading}>More</span>
          <ul className={styles.linkList}>
            <li>
              <a
                className={styles.link}
                href="https://github.com/antoniocali/oat-latte"
                target="_blank"
                rel="noopener noreferrer"
              >
                GitHub&nbsp;&#8599;
              </a>
            </li>
            <li>
              <a
                className={styles.link}
                href="https://github.com/antoniocali/oat-latte/issues"
                target="_blank"
                rel="noopener noreferrer"
              >
                Issues&nbsp;&#8599;
              </a>
            </li>
            <li>
              <a
                className={styles.link}
                href="https://github.com/antoniocali/oat-latte/releases"
                target="_blank"
                rel="noopener noreferrer"
              >
                Releases&nbsp;&#8599;
              </a>
            </li>
          </ul>
        </nav>

      </div>

      {/* Bottom bar */}
      <div className={styles.bottom}>
        <span className={styles.builtWith}>
          Built with&nbsp;
          <a
            className={styles.bottomLink}
            href="https://docusaurus.io"
            target="_blank"
            rel="noopener noreferrer"
          >
            Docusaurus
          </a>
          &nbsp;&mdash; Made with &lt;3 by&nbsp;
          <a
            className={styles.bottomLink}
            href="https://www.linkedin.com/in/antoniodavidecali/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Antonio Davide Cali
          </a>
        </span>
        <span className={styles.mono}>$ go get github.com/antoniocali/oat-latte</span>
      </div>
    </footer>
  );
}
