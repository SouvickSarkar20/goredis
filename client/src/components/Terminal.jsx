import { useEffect, useRef } from 'react'
import ResponseLine from './ResponseLine.jsx'

/**
 * Terminal
 *
 * The main scrollable output area. Renders the startup banner
 * once, then maps each line in `lines` to a ResponseLine component.
 *
 * Auto-scrolls to the bottom whenever new lines are added —
 * same behaviour as a real terminal.
 *
 * Props:
 *   lines  — array of line objects from useRedis
 */
export default function Terminal({ lines }) {
  const bottomRef = useRef(null)

  // Scroll to bottom whenever lines change
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [lines])

  return (
    <div style={styles.terminal}>

      {/* Startup banner — shown once, same style as redis-cli */}
      <Banner />

      {/* All output lines */}
      {lines.map(line => (
        <ResponseLine key={line.id} line={line} />
      ))}

      {/* Invisible anchor element at the bottom for auto-scroll */}
      <div ref={bottomRef} />

    </div>
  )
}

/**
 * Banner
 *
 * The header printed when redis-cli connects. Shows build info
 * and a hint. Purely decorative — makes it feel authentic.
 */
function Banner() {
  return (
    <div style={styles.banner}>
      <div style={styles.bannerLine}>
        <span style={{ color: 'var(--red)' }}>GoRedis</span>
        <span style={{ color: 'var(--comment)' }}> — built from scratch in Go</span>
      </div>
      <div style={styles.bannerLine}>
        <span style={{ color: 'var(--comment)' }}>
          Type a command below or click one from the sidebar.
          Use <span style={{ color: 'var(--yellow)' }}>↑ ↓</span> to recall history.
        </span>
      </div>
      <div style={{ ...styles.bannerLine, marginTop: 4 }}>
        <span style={{ color: 'var(--comment)' }}>Try: </span>
        {['PING', 'SET name Alice', 'GET name', 'LPUSH tasks "write tests"'].map((ex, i) => (
          <span key={i}>
            <span style={styles.exampleCmd}>{ex}</span>
            {i < 3 && <span style={{ color: 'var(--comment)' }}> · </span>}
          </span>
        ))}
      </div>
      <div style={styles.divider} />
    </div>
  )
}

const styles = {
  terminal: {
    flex: 1,
    overflowY: 'auto',
    padding: '14px 20px 12px',
    lineHeight: 1.6,
  },
  banner: {
    marginBottom: 6,
  },
  bannerLine: {
    fontSize: 12,
    marginBottom: 2,
  },
  exampleCmd: {
    color: 'var(--cyan)',
    fontStyle: 'italic',
  },
  divider: {
    borderBottom: '1px solid #292e42',
    marginTop: 10,
    marginBottom: 4,
  },
}
